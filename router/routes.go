package router

import (
	"server_gin/bootstrap"
	"server_gin/handler"
	"server_gin/middleware"

	"github.com/gin-gonic/gin"
)

func registerHealth(r *gin.Engine) {
	h := r.Group("/health")
	Handle(h, "/", handler.Health)
	Handle(h, "/ready", handler.HealthReady)
	Handle(h, "/detail", handler.HealthDetail)
	Handle(h, "/lb", handler.HealthLb)
}

func registerControl(r *gin.Engine) {
	c := r.Group("/internal/control")
	c.Use(middleware.InternalOnly())
	c.POST("/traffic-shift", handler.TrafficShift)
	c.POST("/reject-new-requests", handler.RejectNew)
	c.POST("/resume-traffic", handler.ResumeTraffic)
}

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

func registerOther(group *gin.RouterGroup) {
	other := group.Group("/other")
	other.Use(middleware.DrainGuard())
	Handle(other, "/res_version/:key", handler.ResVersion)
}

func registerTest(group *gin.RouterGroup) {
	test := group.Group("/test")
	test.Use(middleware.TestEnvGuard(), middleware.DrainGuard())
	Handle(test, "/", handler.TestIndex)
	Handle(test, "/db", handler.TestDb)
}
