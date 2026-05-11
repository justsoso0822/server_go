package service

import (
	"context"

	"server_go/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type IUser interface {
	Login(ctx context.Context, uid int64, loginKey, openid, platform, version string) (g.Map, error)
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
