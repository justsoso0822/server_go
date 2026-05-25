package test

import (
	"context"

	v1 "server_go/api/test/v1"
	testService "server_go/internal/service/test"
)

func (c *ControllerV1) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	ret, err := testService.Index(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.IndexRes{
		Data: ret,
	}, nil
}
