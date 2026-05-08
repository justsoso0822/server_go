package service

import (
	"context"

	"server_go/internal/model"
)

type IGrid interface {
	GetGrid(ctx context.Context, in *model.BagInput) (*model.GridOutput, error)
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
