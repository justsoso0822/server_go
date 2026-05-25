package test

import (
	"context"

	v1 "server_go/api/test/v1"
	testService "server_go/internal/service/test"
)

func (c *ControllerV1) TestDb(ctx context.Context, req *v1.TestDbReq) (res *v1.TestDbRes, err error) {
	ret, err := testService.TestDb(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.TestDbRes{Data: ret}, nil
}
