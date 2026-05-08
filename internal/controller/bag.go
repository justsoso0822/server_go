package controller

import (
	"context"

	apiBag "server_go/api/bag"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

var Bag = &cBag{}

type cBag struct{}

func (c *cBag) GetBag(ctx context.Context, req *apiBag.GetBagReq) (res *apiBag.GetBagRes, err error) {
	result, err := service.Bag().GetUserBag(ctx, req.Uid, req.Chapter)
	if err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(result)
	return
}

func (c *cBag) GetBagTp(ctx context.Context, req *apiBag.GetBagTpReq) (res *apiBag.GetBagTpRes, err error) {
	result, err := service.Bag().GetUserBagTp(ctx, req.Uid, req.Chapter)
	if err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(result)
	return
}