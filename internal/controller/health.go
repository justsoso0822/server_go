package controller

import (
	"context"
	"os"
	"sync"
	"time"

	apiHealth "server_go/api/health"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Health = &cHealth{}

type cHealth struct{}

var startTime = time.Now()

var drainState = &drainStateManager{}

type drainStateManager struct {
	mu                   sync.RWMutex
	draining             bool
	rejectingNewRequests bool
}

func (d *drainStateManager) IsTrafficShift() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.draining
}

func (d *drainStateManager) IsRejecting() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.rejectingNewRequests
}

func (d *drainStateManager) StartTrafficShift() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.draining = true
	d.rejectingNewRequests = false
}

func (d *drainStateManager) StartRejectNew() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.draining = true
	d.rejectingNewRequests = true
}

func (d *drainStateManager) Resume() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.draining = false
	d.rejectingNewRequests = false
}

func (c *cHealth) Ready(ctx context.Context, req *apiHealth.ReadyReq) (res *apiHealth.ReadyRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{"ok": true})
	return
}

func (c *cHealth) Health(ctx context.Context, req *apiHealth.HealthReq) (res *apiHealth.HealthRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"status":    "ok",
		"pid":       os.Getpid(),
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": time.Now().Format("2006/01/02 15:04:05"),
	})
	return
}

func (c *cHealth) HealthDetail(ctx context.Context, req *apiHealth.HealthDetailReq) (res *apiHealth.HealthDetailRes, err error) {
	ghttp.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"status":    "ok",
		"pid":       os.Getpid(),
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": time.Now().Format("2006/01/02 15:04:05"),
		"draining":  drainState.IsTrafficShift(),
		"rejecting": drainState.IsRejecting(),
	})
	return
}

func (c *cHealth) HealthLb(ctx context.Context, req *apiHealth.HealthLbReq) (res *apiHealth.HealthLbRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if drainState.IsTrafficShift() {
		r.Response.Status = 503
		r.Response.WriteJson(g.Map{"status": "draining"})
	} else {
		r.Response.WriteJson(g.Map{"status": "ok"})
	}
	return
}

func (c *cHealth) TrafficShift(ctx context.Context, req *apiHealth.TrafficShiftReq) (res *apiHealth.TrafficShiftRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainState.StartTrafficShift()
	r.Response.WriteJson(g.Map{"ok": true, "state": "traffic-shift"})
	return
}

func (c *cHealth) RejectNew(ctx context.Context, req *apiHealth.RejectNewReq) (res *apiHealth.RejectNewRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainState.StartRejectNew()
	r.Response.WriteJson(g.Map{"ok": true, "state": "reject-new-requests"})
	return
}

func (c *cHealth) ResumeTraffic(ctx context.Context, req *apiHealth.ResumeTrafficReq) (res *apiHealth.ResumeTrafficRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if !ensureInternalAccess(r) {
		return
	}
	drainState.Resume()
	r.Response.WriteJson(g.Map{"ok": true, "state": "resume-traffic"})
	return
}

func ensureInternalAccess(r *ghttp.Request) bool {
	expected := os.Getenv("APP_CONTROL_TOKEN")
	if expected == "" || expected == "PLEASE_CHANGE_ME" {
		r.Response.Status = 500
		r.Response.WriteJson(g.Map{"ok": false, "msg": "APP_CONTROL_TOKEN not configured"})
		return false
	}
	forwarded := r.GetHeader("x-forwarded-for")
	if forwarded != "" {
		r.Response.Status = 404
		r.Response.WriteJson(g.Map{"ok": false})
		return false
	}
	received := r.GetHeader("x-control-token")
	if received != expected {
		r.Response.Status = 404
		r.Response.WriteJson(g.Map{"ok": false})
		return false
	}
	return true
}
