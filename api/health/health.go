// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package health

import (
	"context"

	"server_go/api/health/v1"
)

type IHealthV1 interface {
	Ready(ctx context.Context, req *v1.ReadyReq) (res *v1.ReadyRes, err error)
	Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error)
	HealthDetail(ctx context.Context, req *v1.HealthDetailReq) (res *v1.HealthDetailRes, err error)
	HealthLb(ctx context.Context, req *v1.HealthLbReq) (res *v1.HealthLbRes, err error)
}
