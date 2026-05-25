package res

import (
	"context"

	"server_go/api/res/v1"
	resService "server_go/internal/service/res"
)

func (c *ControllerV1) AddGold(ctx context.Context, req *v1.AddGoldReq) (res *v1.AddGoldRes, err error) {
	out, err := resService.UpdateGold(ctx, req.Uid, 50, "测试增加金币")
	if err != nil {
		return nil, err
	}
	return &v1.AddGoldRes{
		Res:     out["res"],
		AddGold: out["add_value"].(int64),
	}, nil
}
