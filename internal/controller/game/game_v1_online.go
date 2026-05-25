package game

import (
	"context"

	"server_go/api/game/v1"
	gameService "server_go/internal/service/game"

	"github.com/gogf/gf/v2/os/gtime"
)

func (c *ControllerV1) Online(ctx context.Context, req *v1.OnlineReq) (res *v1.OnlineRes, err error) {
	err = gameService.Online(ctx, req.Uid, req.Seconds)
	if err != nil {
		return nil, err
	}
	return &v1.OnlineRes{Now: gtime.TimestampMilli()}, nil
}
