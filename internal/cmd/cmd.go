package cmd

import (
	"context"
	"os"
	"strings"

	bagController "server_go/internal/controller/bag"
	controlController "server_go/internal/controller/control"
	gameController "server_go/internal/controller/game"
	gridController "server_go/internal/controller/grid"
	healthController "server_go/internal/controller/health"
	otherController "server_go/internal/controller/other"
	resController "server_go/internal/controller/res"
	testController "server_go/internal/controller/test"
	userController "server_go/internal/controller/user"
	"server_go/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcmd"

	"server_go/utility/dbcache"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 根据配置决定是否启用 ORM 查询缓存
			if g.Cfg().MustGet(ctx, "database.default.cache").Bool() {
				g.DB().GetCache().SetAdapter(gcache.NewAdapterRedis(g.Redis()))
			} else {
				g.DB().GetCache().SetAdapter(&dbcache.NoopAdapter{})
			}

			s := g.Server()
			if appPort := strings.TrimSpace(os.Getenv("APP_PORT")); appPort != "" {
				s.SetAddr(":" + strings.TrimPrefix(appPort, ":"))
			}

			// 内部控制路由不走渠道分流，供部署脚本直接控制当前容器状态。
			s.Group("/internal/control", func(group *ghttp.RouterGroup) {
				group.Bind(controlController.NewV1())
			})

			// 健康检查路由不走渠道分流，供 Traefik 和部署脚本使用固定路径。
			s.Group("/health", func(group *ghttp.RouterGroup) {
				group.Bind(healthController.NewV1())
			})

			registerChannelRoutes := func(channel *ghttp.RouterGroup) {
				// 游戏接口路由
				channel.Group("/api", func(group *ghttp.RouterGroup) {
					group.Middleware(
						middleware.DrainGuard,
						middleware.Sign,
						middleware.Verify,
						ghttp.MiddlewareHandlerResponse,
					)
					group.Bind(
						userController.NewV1(),
						gameController.NewV1(),
						bagController.NewV1(),
						gridController.NewV1(),
						resController.NewV1(),
					)
				})

				// 其他路由（不校验签名和登录态）
				channel.Group("/other", func(group *ghttp.RouterGroup) {
					group.Middleware(
						middleware.DrainGuard,
						ghttp.MiddlewareHandlerResponse,
					)
					group.Bind(otherController.NewV1())
				})

				// 测试路由（仅 local/test 环境可用）
				channel.Group("/test", func(group *ghttp.RouterGroup) {
					group.Middleware(
						middleware.TestEnvGuard,
						middleware.DrainGuard,
						ghttp.MiddlewareHandlerResponse,
					)
					group.Bind(testController.NewV1())
				})
			}

			// 无渠道前缀入口默认走 default 渠道。
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.DefaultChannel)
				registerChannelRoutes(group)
			})

			// 带渠道前缀入口走指定渠道。
			s.Group("/{channel}", func(channel *ghttp.RouterGroup) {
				channel.Middleware(middleware.Channel)
				registerChannelRoutes(channel)
			})

			s.Run()
			return nil
		},
	}
)
