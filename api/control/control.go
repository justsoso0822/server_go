// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package control

import (
	"context"

	"server_go/api/control/v1"
)

type IControlV1 interface {
	TrafficShift(ctx context.Context, req *v1.TrafficShiftReq) (res *v1.TrafficShiftRes, err error)
	RejectNew(ctx context.Context, req *v1.RejectNewReq) (res *v1.RejectNewRes, err error)
	ResumeTraffic(ctx context.Context, req *v1.ResumeTrafficReq) (res *v1.ResumeTrafficRes, err error)
}
