package bag

import (
	"context"

	"server_go/api/bag/v1"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
)

func (c *ControllerV1) GetBagTp(ctx context.Context, req *v1.GetBagTpReq) (res *v1.GetBagTpRes, err error) {
	out, err := service.Bag().GetUserBagTp(ctx, req.Uid, req.Chapter)
	if err != nil {
		return nil, err
	}
	return &v1.GetBagTpRes{
		Uid:     out["uid"].(int64),
		Chapter: out["chapter"].(int),
		Bag:     out["bag"].(gdb.Result),
	}, nil
}
