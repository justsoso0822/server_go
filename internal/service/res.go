package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	IRes interface {
		UpdateDiamond(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error)
		UpdateGold(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error)
		UpdateTili(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error)
		UpdateExp(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error)
		UpdateStar(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error)
	}
)

var (
	localRes IRes
)

func Res() IRes {
	if localRes == nil {
		panic("implement not found for interface IRes, forgot register?")
	}
	return localRes
}

func RegisterRes(i IRes) {
	localRes = i
}
