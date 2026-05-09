package main

import (
	_ "server_go/internal/packed"

	_ "server_go/internal/logic/bag"
	_ "server_go/internal/logic/game"
	_ "server_go/internal/logic/grid"
	_ "server_go/internal/logic/other"
	_ "server_go/internal/logic/res"
	_ "server_go/internal/logic/task"
	_ "server_go/internal/logic/user"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	_ "server_go/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"server_go/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
