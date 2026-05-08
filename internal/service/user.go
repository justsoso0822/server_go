package service

import "context"

type IUser interface {
	Login(ctx context.Context, uid int64, loginKey, openid, platform, version string) (map[string]interface{}, error)
	UpdateDiamond(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error
	UpdateGold(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error
	UpdateTili(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error
	UpdateExp(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error
	UpdateStar(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error
	UpdateRes(ctx context.Context, uid int64, items interface{}, ret map[string]interface{}, reason string) error
	UpdateItemsOther(ctx context.Context, uid int64, items map[int]int, ret map[string]interface{}, reason string) error
	GetUser(ctx context.Context, uid int64) (map[string]interface{}, error)
	GetUserRes(ctx context.Context, uid int64) (map[string]interface{}, error)
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