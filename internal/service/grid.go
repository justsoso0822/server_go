package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type IGrid interface {
	GetGrid(ctx context.Context, uid int64, chapter int) (g.Map, error)
}

var localGrid IGrid

func Grid() IGrid {
	if localGrid == nil {
		panic("service IGrid not registered")
	}
	return localGrid
}

func RegisterGrid(s IGrid) {
	localGrid = s
}
