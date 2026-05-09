// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package bag

import (
	"context"

	"server_go/api/bag/v1"
)

type IBagV1 interface {
	GetBag(ctx context.Context, req *v1.GetBagReq) (res *v1.GetBagRes, err error)
	GetBagTp(ctx context.Context, req *v1.GetBagTpReq) (res *v1.GetBagTpRes, err error)
}
