// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	IGrid interface {
		GetGrid(ctx context.Context, uid int64, chapter int) (g.Map, error)
	}
)

var (
	localGrid IGrid
)

func Grid() IGrid {
	if localGrid == nil {
		panic("implement not found for interface IGrid, forgot register?")
	}
	return localGrid
}

func RegisterGrid(i IGrid) {
	localGrid = i
}
