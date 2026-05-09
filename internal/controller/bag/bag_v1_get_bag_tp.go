package bag

import (
	"context"

	"server_go/api/bag/v1"
	"server_go/internal/model"
	"server_go/internal/service"
)

func (c *ControllerV1) GetBagTp(ctx context.Context, req *v1.GetBagTpReq) (res *v1.GetBagTpRes, err error) {
	out, err := service.Bag().GetUserBagTp(ctx, &model.BagInput{Uid: req.Uid, Chapter: req.Chapter})
	if err != nil {
		return nil, err
	}
	return (*v1.GetBagTpRes)(ToBagRes(out)), nil
}
