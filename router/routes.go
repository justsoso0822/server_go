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

// 正式业务接口统一挂在 /api 下，并按固定顺序经过中间件链。
// Gin 会按 Use 的注册顺序执行：先处理排水，再验签，再防重放，最后校验登录态。
// 这个顺序能尽早拒绝无效请求，减少 DB/Redis 热路径之外的业务开销。
func registerAPI(group *gin.RouterGroup, app *bootstrap.App) {
	api := group.Group("/api")
	api.Use(
		middleware.DrainGuard(),
		middleware.Sign(app.Config),
		middleware.ReplayGuard(),
		middleware.Verify(),
	)
	Handle(api, "/user/login", handler.UserLogin)
	Handle(api, "/game/time", handler.GameTime)
	Handle(api, "/game/online", handler.GameOnline)
	Handle(api, "/bag/get_bag/:chapter", handler.GetBag)
	Handle(api, "/bag/get_bag_tp/:chapter", handler.GetBagTp)
	Handle(api, "/grid/get/:chapter", handler.GetGrid)
	Handle(api, "/res/add_tili", handler.AddTili)
	Handle(api, "/res/add_gold", handler.AddGold)
	Handle(api, "/res/add_diamond", handler.AddDiamond)
}
