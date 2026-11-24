package router

import (
	"io"
	"os"

	"app/bootstrap"
	"app/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *bootstrap.App) *gin.Engine {
	configureGinMode(app.Config.App.Env)

	// 使用 gin.New 而不是 gin.Default，是为了不启用 Gin 自带 Logger/Recovery。
	// 本项目自行接入 zap：日志格式统一，panic 栈也能带 request_id/channel 等字段。
	r := gin.New()
	r.Use(middleware.RequestID())
	// CORS 必须排在验签、渠道等业务中间件之前：
	// 浏览器预检 OPTIONS 不带 sign/channel 参数，先进入业务链会被误判为非法调用。
	r.Use(middleware.CORS(app.Config.App.AllowedOrigins, app.Config.App.Env))
	// 请求级超时：给下游 Redis / MySQL 一个统一的截止时间。
	// 正常请求远低于 2 秒，故障时也能在「HTTP 层自己超时」之前快速失败。
	r.Use(middleware.RequestTimeout(middleware.DefaultRequestTimeout))
	r.Use(middleware.AccessLog(app.Logger))
	r.Use(middleware.Recovery(app.Logger))

	registerHealth(r)
	registerControl(r)

	defaultGroup := r.Group("/")
	defaultGroup.Use(middleware.DefaultChannel())
	bindChannelRoutes(defaultGroup, app)

	channelGroup := r.Group("/:channel")
	// 带 /:channel 前缀时，Channel 中间件把路由参数写进 request context；
	// 后续 autodb.DB/Redis 会据此选择对应渠道的 MySQL/Redis 连接。
	channelGroup.Use(middleware.Channel())
	bindChannelRoutes(channelGroup, app)

	return r
}

func bindChannelRoutes(group *gin.RouterGroup, app *bootstrap.App) {
	registerOther(group)
	registerTest(group, app.Config.App.Env)
	registerAPI(group, app)
}

func configureGinMode(env string) {
	if env == "" || env == "local" {
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
		return
	}

	gin.SetMode(gin.ReleaseMode)
	// Release 模式下丢弃 Gin 默认输出，避免框架裸文本日志混进 zap JSON 日志流。
	// 真正的访问日志和 panic 日志都在自定义中间件里输出。
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}
