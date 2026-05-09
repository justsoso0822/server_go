package res

import (
	"context"

	"server_go/api/res/v1"
	"server_go/internal/model"
	"server_go/internal/service"
)

func (c *ControllerV1) AddGold(ctx context.Context, req *v1.AddGoldReq) (res *v1.AddGoldRes, err error) {
	out, err := service.Res().UpdateGold(ctx, &model.UpdateFieldInput{Uid: req.Uid, Cnt: 50, Reason: "测试增加金币"})
	if err != nil {
		return nil, err
	}
	return &v1.AddGoldRes{
		Res:     out.Res,
		AddGold: out.AddValue,
	}, nil
}
