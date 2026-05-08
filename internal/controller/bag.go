package controller

import (
	"context"

	apiBag "server_go/api/bag"
	"server_go/internal/model"
	"server_go/internal/service"
)

var Bag = &cBag{}

type cBag struct{}

func (c *cBag) GetBag(ctx context.Context, req *apiBag.GetBagReq) (res *apiBag.GetBagRes, err error) {
	out, err := service.Bag().GetUserBag(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return (*apiBag.GetBagRes)(toBagRes(out)), nil
}

func (c *cBag) GetBagTp(ctx context.Context, req *apiBag.GetBagTpReq) (res *apiBag.GetBagTpRes, err error) {
	out, err := service.Bag().GetUserBagTp(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return (*apiBag.GetBagTpRes)(toBagRes(out)), nil
}

func toBagRes(out *model.BagOutput) *apiBag.BagRes {
	if out == nil {
		return nil
	}
	return &apiBag.BagRes{
		Uid:     out.Uid,
		Chapter: out.Chapter,
		Bag:     out.Bag,
	}
}
