package cmd

import (
	"context"
	"os"
	"strings"

	"server_go/internal/autodb"
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
	"server_go/utility/dbcache"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "start http server",
	Func:  run,
}

func run(ctx context.Context, _ *gcmd.Parser) error {
	channels, err := autodb.LoadConfiguredChannels(ctx)
	if err != nil {
		g.Log().Fatalf(ctx, "load channels: %v", err)
		return err
	}
	g.Log().Infof(ctx, "configured channels: %v", channels)

	for _, group := range channels {
		applyDBCache(ctx, group)
	}

	s := g.Server()
	applyServerAddr(s)
	registerSystemRoutes(s)
	registerChannelEntries(s)
	s.Run()
	return nil
}

// applyDBCache 按 database.<group>.cache 决定该分组的 ORM 查询缓存适配器。
// 启用时以同 channel 的 Redis 作为载体，否则使用空操作适配器。
func applyDBCache(ctx context.Context, group string) {
	if g.Cfg().MustGet(ctx, "database."+group+".cache").Bool() {
		g.DB(group).GetCache().SetAdapter(gcache.NewAdapterRedis(g.Redis(group)))
		return
	}
	g.DB(group).GetCache().SetAdapter(&dbcache.NoopAdapter{})
}

// applyServerAddr 在 APP_PORT 环境变量存在时覆盖配置中的监听端口。
func applyServerAddr(s *ghttp.Server) {
	appPort := strings.TrimSpace(os.Getenv("APP_PORT"))
	if appPort == "" {
		return
	}
	s.SetAddr(":" + strings.TrimPrefix(appPort, ":"))
}

// registerSystemRoutes 注册不参与渠道分流的系统级路由。
// 内部控制路由供部署脚本直接控制容器状态，健康检查路由供 Traefik 使用固定路径。
func registerSystemRoutes(s *ghttp.Server) {
	s.Group("/internal/control", func(group *ghttp.RouterGroup) {
		group.Bind(controlController.NewV1())
	})
	s.Group("/health", func(group *ghttp.RouterGroup) {
		group.Bind(healthController.NewV1())
	})
}

// registerChannelEntries 注册两个 channel 入口：
//   - "/"             默认走 default 渠道
//   - "/{channel}"    路径首段为 channel 名
func registerChannelEntries(s *ghttp.Server) {
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.DefaultChannel)
		bindChannelRoutes(group)
	})
	s.Group("/{channel}", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Channel)
		bindChannelRoutes(group)
	})
}

// bindChannelRoutes 在每个 channel 入口下注册业务子路径与对应中间件。
func bindChannelRoutes(channel *ghttp.RouterGroup) {
	// 游戏接口：校验签名 + 登录态
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

	// 其他接口：不校验签名和登录态
	channel.Group("/other", func(group *ghttp.RouterGroup) {
		group.Middleware(
			middleware.DrainGuard,
			ghttp.MiddlewareHandlerResponse,
		)
		group.Bind(otherController.NewV1())
	})

	// 测试接口：仅 local/test 环境可用
	channel.Group("/test", func(group *ghttp.RouterGroup) {
		group.Middleware(
			middleware.TestEnvGuard,
			middleware.DrainGuard,
			ghttp.MiddlewareHandlerResponse,
		)
		group.Bind(testController.NewV1())
	})
}