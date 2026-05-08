package controller

import (
	"context"

	apiOther "server_go/api/other"
	"server_go/internal/service"
)

var Other = &cOther{}

type cOther struct{}

func (c *cOther) ResVersion(ctx context.Context, req *apiOther.ResVersionReq) (res *apiOther.ResVersionRes, err error) {
	out, err := service.Other().GetResVersion(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &apiOther.ResVersionRes{
		Code: out.Code,
		Ver:  out.Ver,
		Msg:  out.Msg,
	}, nil
}
