// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package test

import (
	"context"

	"server_go/api/test/v1"
)

type ITestV1 interface {
	Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error)
}
