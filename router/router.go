package router

import (
	"server_gin/bootstrap"
	"server_gin/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *bootstrap.App) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	registerHealth(r)
	registerControl(r)

	defaultGroup := r.Group("/")
	defaultGroup.Use(middleware.DefaultChannel())
	bindChannelRoutes(defaultGroup, app)

	channelGroup := r.Group("/:channel")
	channelGroup.Use(middleware.Channel())
	bindChannelRoutes(channelGroup, app)

	return r
}

func bindChannelRoutes(group *gin.RouterGroup, app *bootstrap.App) {
	registerAPI(group, app)
	registerOther(group)
	registerTest(group)
}
