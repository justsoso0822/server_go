package controller

import (
	"context"

	apiGame "server_go/api/game"
	"server_go/internal/model"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
)

var Game = &cGame{}

type cGame struct{}

func (c *cGame) Online(ctx context.Context, req *apiGame.OnlineReq) (res *apiGame.OnlineRes, err error) {
	err = service.Game().Online(ctx, &model.OnlineInput{
		Uid: req.Uid, Seconds: req.Seconds,
	})
	if err != nil {
		return nil, err
	}
	return &apiGame.OnlineRes{Now: gtime.TimestampMilli()}, nil
}

func (c *cGame) Time(ctx context.Context, req *apiGame.TimeReq) (res *apiGame.TimeRes, err error) {
	return &apiGame.TimeRes{Now: gtime.TimestampMilli()}, nil
}
