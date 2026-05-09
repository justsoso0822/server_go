package game

import (
	"context"

	"server_go/api/game/v1"

	"github.com/gogf/gf/v2/os/gtime"
)

func (c *ControllerV1) Time(ctx context.Context, req *v1.TimeReq) (res *v1.TimeRes, err error) {
	return &v1.TimeRes{Now: gtime.TimestampMilli()}, nil
}
