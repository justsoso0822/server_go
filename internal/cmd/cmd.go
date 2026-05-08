package cmd

import (
	"context"

	"server_go/internal/controller"
	"server_go/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// Game API routes
			s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(
					middleware.Sign,
					middleware.Verify,
					middleware.Response,
				)
				group.Bind(
					controller.User,
					controller.Game,
					controller.Bag,
					controller.Grid,
				)
			})

			// Other routes (no sign/verify)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Response)
				group.Bind(
					controller.Other,
				)
			})

			// Health + internal control routes (no middleware)
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Bind(
					controller.Health,
				)
			})

			s.Run()
			return nil
		},
	}
)
