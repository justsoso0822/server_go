package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type TrafficShiftReq struct {
	g.Meta `path:"/traffic-shift" method:"get,post" tags:"Internal" summary:"开始流量切换"`
}
type TrafficShiftRes struct {
	g.Meta `mime:"application/json"`
}

type RejectNewReq struct {
	g.Meta `path:"/reject-new-requests" method:"get,post" tags:"Internal" summary:"拒绝新请求"`
}
type RejectNewRes struct {
	g.Meta `mime:"application/json"`
}

type ResumeTrafficReq struct {
	g.Meta `path:"/resume-traffic" method:"get,post" tags:"Internal" summary:"恢复流量"`
}
type ResumeTrafficRes struct {
	g.Meta `mime:"application/json"`
}
