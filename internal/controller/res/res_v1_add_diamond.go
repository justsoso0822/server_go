package res

import (
	"context"

	"server_go/api/res/v1"
	resService "server_go/internal/service/res"
)

func (c *ControllerV1) AddDiamond(ctx context.Context, req *v1.AddDiamondReq) (res *v1.AddDiamondRes, err error) {
	out, err := resService.UpdateDiamond(ctx, req.Uid, 50, "测试增加钻石")
	if err != nil {
		return nil, err
	}
	return &v1.AddDiamondRes{
		Res:        out["res"],
		AddDiamond: out["add_value"].(int64),
	}, nil
}
