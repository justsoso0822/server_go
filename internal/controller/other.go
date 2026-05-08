package controller

import (
	"context"

	apiOther "server_go/api/other"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

var Other = &cOther{}

type cOther struct{}

func (c *cOther) ResVersion(ctx context.Context, req *apiOther.ResVersionReq) (res *apiOther.ResVersionRes, err error) {
	result, err := service.Other().GetResVersion(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(result)
	return
}