package service

import (
	"context"

	"server_go/internal/model"
	"server_go/internal/model/entity"
)

type IUser interface {
	Login(ctx context.Context, in *model.LoginInput) (*model.LoginOutput, error)
	UpdateDiamond(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateGold(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateTili(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateExp(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	UpdateStar(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error)
	GetUser(ctx context.Context, uid int64) (*entity.User, error)
	GetUserRes(ctx context.Context, uid int64) (*entity.UserRes, error)
}

var localUser IUser

func User() IUser {
	if localUser == nil {
		panic("service IUser not registered")
	}
	return localUser
}

func RegisterUser(s IUser) {
	localUser = s
}