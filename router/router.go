package router

import (
	"io"
	"os"
	"strings"

	"server_gin/bootstrap"
	"server_gin/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *bootstrap.App) *gin.Engine {
	configureGinMode()

	r := gin.New()
	r.Use(middleware.Recovery(app.Logger))

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

func configureGinMode() {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env == "" || env == "local" {
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = os.Stdout
		gin.DefaultErrorWriter = os.Stderr
		return
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}
