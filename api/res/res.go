// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package res

import (
	"context"

	"server_go/api/res/v1"
)

type IResV1 interface {
	AddTili(ctx context.Context, req *v1.AddTiliReq) (res *v1.AddTiliRes, err error)
	AddGold(ctx context.Context, req *v1.AddGoldReq) (res *v1.AddGoldRes, err error)
	AddDiamond(ctx context.Context, req *v1.AddDiamondReq) (res *v1.AddDiamondRes, err error)
}
