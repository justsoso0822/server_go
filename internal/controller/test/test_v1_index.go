package test

import (
	"context"

	v1 "server_go/api/test/v1"
	"server_go/internal/service"
)

func (c *ControllerV1) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	ret, err := service.Test().Index(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.IndexRes{
		Data: ret,
	}, nil
}
