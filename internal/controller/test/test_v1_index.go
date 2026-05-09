package test

import (
	"context"

	"server_go/api/test/v1"
)

func (c *ControllerV1) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	return &v1.IndexRes{Data: "ok"}, nil
}
