package test

import (
	"context"

	v1 "server_go/api/test/v1"
	"server_go/internal/service"
)

func (c *ControllerV1) TestDb(ctx context.Context, req *v1.TestDbReq) (res *v1.TestDbRes, err error) {
	ret, err := service.Test().TestDb(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.TestDbRes{Data: ret}, nil
}
