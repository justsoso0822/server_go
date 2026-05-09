package grid

import (
	"context"

	"server_go/api/grid/v1"
	bagController "server_go/internal/controller/bag"
	"server_go/internal/model"
	"server_go/internal/service"
)

func (c *ControllerV1) GetGrid(ctx context.Context, req *v1.GetGridReq) (res *v1.GetGridRes, err error) {
	out, err := service.Grid().GetGrid(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return &v1.GetGridRes{
		Bag:   bagController.ToBagRes(out.Bag),
		BagTp: bagController.ToBagRes(out.BagTp),
		Tasks: out.Tasks,
	}, nil
}
