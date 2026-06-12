package router

import (
	"server_go/bootstrap"
	"server_go/handler"
	"server_go/middleware"

	"github.com/gin-gonic/gin"
)

// 健康检查接口挂在根级路由上，不走渠道、签名和登录校验。
// 这些接口通常给容器 healthcheck、Traefik 负载均衡和人工排查使用。
func registerHealth(r *gin.Engine) {
	h := r.Group("/health")
	Handle(h, "/", handler.Health)
	Handle(h, "/ready", handler.HealthReady)
	Handle(h, "/detail", handler.HealthDetail)
	Handle(h, "/lb", handler.HealthLb)
}

// 内部控制接口只允许容器内网直连访问，用于蓝绿部署期间切流和排水。
// 这里使用 Gin 的 RouterGroup.Use 给整组接口加 middleware.InternalOnly，
// 避免每个 POST 路由重复写同一段访问控制逻辑。
func registerControl(r *gin.Engine) {
	c := r.Group("/internal/control")
	c.Use(middleware.InternalOnly())
	c.POST("/traffic-shift", handler.TrafficShift)
	c.POST("/reject-new-requests", handler.RejectNew)
	c.POST("/resume-traffic", handler.ResumeTraffic)
}

// 其他辅助接口依赖渠道上下文，但不进入正式 API 的签名和登录校验链。
// DrainGuard 会在服务排水阶段拒绝新请求，防止旧实例继续接收业务流量。
func registerOther(group *gin.RouterGroup) {
	other := group.Group("/other")
	other.Use(middleware.DrainGuard())
	Handle(other, "/res_version/:key", handler.ResVersion)
}

// 测试接口只在 local/test 环境开放，并且同样受排水状态控制。
// TestEnvGuard 放在前面可以先挡住生产访问，再进入业务处理。
func registerTest(group *gin.RouterGroup) {
	test := group.Group("/test")
	test.Use(middleware.TestEnvGuard(), middleware.DrainGuard())
	Handle(test, "/", handler.TestIndex)
	Handle(test, "/db", handler.TestDB)
}

// 正式业务接口统一挂在 /api 下。
// 公共链路只保留排水和验签；防重放只给登录/写接口使用，避免读接口每次请求都写 Redis。
func registerAPI(group *gin.RouterGroup, app *bootstrap.App) {
	api := group.Group("/api")
	api.Use(
		middleware.DrainGuard(),
		middleware.Sign(app.Config),
	)

	// 登录接口还没有登录态，只做签名和防重放。
	login := api.Group("")
	login.Use(middleware.ReplayGuard())
	Handle(login, "/user/login", handler.UserLogin)

	// 读接口只需要签名和登录态校验，不使用 ReplayGuard，减少 Redis SET NX 写压力。
	read := api.Group("")
	read.Use(middleware.Verify())
	Handle(read, "/game/time", handler.GameTime)
	Handle(read, "/game/online", handler.GameOnline)
	Handle(read, "/bag/get_bag/:chapter", handler.GetBag)
	Handle(read, "/bag/get_bag_tp/:chapter", handler.GetBagTp)
	Handle(read, "/grid/get/:chapter", handler.GetGrid)

	// 写接口先校验登录态，再做防重放；无效 login_key 不会消耗 replay key。
	write := api.Group("")
	write.Use(
		middleware.Verify(),
		middleware.ReplayGuard(),
	)
	Handle(write, "/res/add_tili", handler.AddTili)
	Handle(write, "/res/add_gold", handler.AddGold)
	Handle(write, "/res/add_diamond", handler.AddDiamond)
}
