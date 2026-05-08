package controller

import (
	"context"

	apiGrid "server_go/api/grid"
	"server_go/internal/model"
	"server_go/internal/service"
)

var Grid = &cGrid{}

type cGrid struct{}

func (c *cGrid) GetGrid(ctx context.Context, req *apiGrid.GetGridReq) (res *apiGrid.GetGridRes, err error) {
	out, err := service.Grid().GetGrid(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return &apiGrid.GetGridRes{
		Bag:   toBagRes(out.Bag),
		BagTp: toBagRes(out.BagTp),
		Tasks: out.Tasks,
	}, nil
}
