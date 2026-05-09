// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server_go/internal/model"
)

type (
	IRes interface {
		UpdateDiamond(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
		UpdateGold(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
		UpdateTili(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
		UpdateExp(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
		UpdateStar(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
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
