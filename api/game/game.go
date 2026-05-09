// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package game

import (
	"context"

	"server_go/api/game/v1"
)

type IGameV1 interface {
	Online(ctx context.Context, req *v1.OnlineReq) (res *v1.OnlineRes, err error)
	Time(ctx context.Context, req *v1.TimeReq) (res *v1.TimeRes, err error)
}
