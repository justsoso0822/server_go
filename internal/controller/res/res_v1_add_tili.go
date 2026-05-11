package res

import (
	"context"

	"server_go/api/res/v1"
	"server_go/internal/service"
)

func (c *ControllerV1) AddTili(ctx context.Context, req *v1.AddTiliReq) (res *v1.AddTiliRes, err error) {
	out, err := service.Res().UpdateTili(ctx, req.Uid, 50, "测试增加体力")
	if err != nil {
		return nil, err
	}
	return &v1.AddTiliRes{
		Res:     out["res"],
		AddTili: out["add_value"].(int64),
	}, nil
}
