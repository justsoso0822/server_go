package controller

import (
	"context"

	apiGrid "server_go/api/grid"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

var Grid = &cGrid{}

type cGrid struct{}

func (c *cGrid) GetGrid(ctx context.Context, req *apiGrid.GetGridReq) (res *apiGrid.GetGridRes, err error) {
	result, err := service.Grid().GetGrid(ctx, req.Uid, req.Chapter)
	if err != nil {
		return nil, err
	}
	ghttp.RequestFromCtx(ctx).Response.WriteJson(result)
	return
}