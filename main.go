package main

import (
	_ "server_go/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"server_go/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
