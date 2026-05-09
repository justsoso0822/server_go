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

			// 游戏接口路由
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

			// 其他路由（不校验签名和登录态）
			s.Group("/other", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Response)
				group.Bind(
					controller.Other,
				)
			})

			// 健康检查路由（无中间件）
			s.Group("/health", func(group *ghttp.RouterGroup) {
				group.Bind(
					controller.Health,
				)
			})

			// 内部控制路由（无中间件）
			s.Group("/_internal/control", func(group *ghttp.RouterGroup) {
				group.Bind(
					controller.InternalControl,
				)
			})

			s.Run()
			return nil
		},
	}
)
