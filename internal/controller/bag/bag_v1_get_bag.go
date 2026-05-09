package bag

import (
	"context"

	"server_go/api/bag/v1"
	"server_go/internal/model"
	"server_go/internal/service"
)

func (c *ControllerV1) GetBag(ctx context.Context, req *v1.GetBagReq) (res *v1.GetBagRes, err error) {
	out, err := service.Bag().GetUserBag(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return (*v1.GetBagRes)(ToBagRes(out)), nil
}
