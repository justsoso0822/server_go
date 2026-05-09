package health

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ReadyReq struct {
	g.Meta `path:"/ready" method:"get,post" tags:"Health" summary:"就绪检查"`
}
type ReadyRes struct {
	g.Meta `mime:"application/json"`
}

type HealthReq struct {
	g.Meta `path:"/" method:"get,post" tags:"Health" summary:"健康检查"`
}
type HealthRes struct {
	g.Meta `mime:"application/json"`
}

type HealthDetailReq struct {
	g.Meta `path:"/detail" method:"get,post" tags:"Health" summary:"健康详情"`
}
type HealthDetailRes struct {
	g.Meta `mime:"application/json"`
}

type HealthLbReq struct {
	g.Meta `path:"/lb" method:"get,post" tags:"Health" summary:"LB健康检查"`
}
type HealthLbRes struct {
	g.Meta `mime:"application/json"`
}

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
