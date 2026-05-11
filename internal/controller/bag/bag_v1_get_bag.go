package bag

import (
	"context"

	"server_go/api/bag/v1"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
)

func (c *ControllerV1) GetBag(ctx context.Context, req *v1.GetBagReq) (res *v1.GetBagRes, err error) {
	out, err := service.Bag().GetUserBag(ctx, req.Uid, req.Chapter)
	if err != nil {
		return nil, err
	}
	return &v1.GetBagRes{
		Uid:     out["uid"].(int64),
		Chapter: out["chapter"].(int),
		Bag:     out["bag"].(gdb.Result),
	}, nil
}
