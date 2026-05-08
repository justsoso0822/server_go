package controller

import (
	"context"

	apiUser "server_go/api/user"
	"server_go/internal/model"
	"server_go/internal/service"
)

var User = &cUser{}

type cUser struct{}

func (c *cUser) Login(ctx context.Context, req *apiUser.LoginReq) (res *apiUser.LoginRes, err error) {
	out, err := service.User().Login(ctx, &model.LoginInput{
		Uid:      req.Uid,
		LoginKey: req.LoginKey,
		Openid:   req.Openid,
		Platform: req.Platform,
		Version:  req.Version,
	})
	if err != nil {
		return nil, err
	}
	return (*apiUser.LoginRes)(out), nil
}

func (c *cUser) AddTili(ctx context.Context, req *apiUser.AddTiliReq) (res *apiUser.AddTiliRes, err error) {
	out, err := service.User().UpdateTili(ctx, &model.UpdateFieldInput{
		Uid: req.Uid, Cnt: 50, Reason: "测试增加体力",
	})
	if err != nil {
		return nil, err
	}
	return (*apiUser.AddTiliRes)(out), nil
}

func (c *cUser) AddGold(ctx context.Context, req *apiUser.AddGoldReq) (res *apiUser.AddGoldRes, err error) {
	out, err := service.User().UpdateGold(ctx, &model.UpdateFieldInput{
		Uid: req.Uid, Cnt: 50, Reason: "测试增加金币",
	})
	if err != nil {
		return nil, err
	}
	return (*apiUser.AddGoldRes)(out), nil
}

func (c *cUser) AddDiamond(ctx context.Context, req *apiUser.AddDiamondReq) (res *apiUser.AddDiamondRes, err error) {
	out, err := service.User().UpdateDiamond(ctx, &model.UpdateFieldInput{
		Uid: req.Uid, Cnt: 50, Reason: "测试增加钻石",
	})
	if err != nil {
		return nil, err
	}
	return (*apiUser.AddDiamondRes)(out), nil
}