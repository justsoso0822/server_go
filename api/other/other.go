// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package other

import (
	"context"

	"server_go/api/other/v1"
)

type IOtherV1 interface {
	ResVersion(ctx context.Context, req *v1.ResVersionReq) (res *v1.ResVersionRes, err error)
}
