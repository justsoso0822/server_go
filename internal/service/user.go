// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"server_go/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	IUser interface {
		Login(ctx context.Context, uid int64, loginKey string, openid string, platform string, version string) (g.Map, error)
		GetUser(ctx context.Context, uid int64) (*entity.User, error)
		GetUserRes(ctx context.Context, uid int64) (*entity.UserRes, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
