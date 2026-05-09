// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"server_go/api/user/v1"
)

type IUserV1 interface {
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	AddTili(ctx context.Context, req *v1.AddTiliReq) (res *v1.AddTiliRes, err error)
	AddGold(ctx context.Context, req *v1.AddGoldReq) (res *v1.AddGoldRes, err error)
	AddDiamond(ctx context.Context, req *v1.AddDiamondReq) (res *v1.AddDiamondRes, err error)
}
