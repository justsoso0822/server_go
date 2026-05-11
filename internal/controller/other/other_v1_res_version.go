package other

import (
	"context"

	"server_go/api/other/v1"
	"server_go/internal/service"
)

func (c *ControllerV1) ResVersion(ctx context.Context, req *v1.ResVersionReq) (res *v1.ResVersionRes, err error) {
	out, err := service.Other().GetResVersion(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &v1.ResVersionRes{
		Code: out["code"].(int),
		Ver:  out["ver"].(string),
		Msg:  out["msg"].(string),
	}, nil
}
