package router

import (
	"io"
	"os"
	"strings"

	"server_go/bootstrap"
	"server_go/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(app *bootstrap.App) *gin.Engine {
	configureGinMode()

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.AccessLog(app.Logger))
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
	registerOther(group)
	registerTest(group)
	registerAPI(group, app)
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
