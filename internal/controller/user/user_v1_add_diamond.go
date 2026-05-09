package user

import (
	"context"

	"server_go/api/user/v1"
	"server_go/internal/model"
	"server_go/internal/service"
)

func (c *ControllerV1) AddDiamond(ctx context.Context, req *v1.AddDiamondReq) (res *v1.AddDiamondRes, err error) {
	out, err := service.User().UpdateDiamond(ctx, &model.UpdateFieldInput{Uid: req.Uid, Cnt: 50, Reason: "测试增加钻石"})
	if err != nil {
		return nil, err
	}
	return (*v1.AddDiamondRes)(toUpdateFieldRes(out)), nil
}
