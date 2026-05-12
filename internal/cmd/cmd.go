package cmd

import (
	"context"

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

			// 游戏接口路由
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(
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
			s.Group("/other", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					otherController.NewV1(),
				)
			})

			// 健康检查路由（无中间件）
			s.Group("/health", func(group *ghttp.RouterGroup) {
				group.Bind(
					healthController.NewV1(),
				)
			})

			// 内部控制路由（无中间件）
			s.Group("/_internal/control", func(group *ghttp.RouterGroup) {
				group.Bind(
					controlController.NewV1(),
				)
			})

			s.Group("/test", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					testController.NewV1(),
				)
			})

			s.Run()
			return nil
		},
	}
)
