// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package grid

import (
	"context"

	"server_go/api/grid/v1"
)

type IGridV1 interface {
	GetGrid(ctx context.Context, req *v1.GetGridReq) (res *v1.GetGridRes, err error)
}
