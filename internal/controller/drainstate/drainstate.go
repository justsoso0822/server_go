package drainstate

import "github.com/gogf/gf/v2/container/gtype"

var manager = &stateManager{
	draining:             gtype.NewBool(),
	rejectingNewRequests: gtype.NewBool(),
	activeRequests:       gtype.NewInt64(),
}

type stateManager struct {
	draining             *gtype.Bool
	rejectingNewRequests *gtype.Bool
	activeRequests       *gtype.Int64
}

func IncActiveRequests() {
	manager.activeRequests.Add(1)
}

func DecActiveRequests() {
	manager.activeRequests.Add(-1)
}

func GetActiveRequests() int64 {
	return manager.activeRequests.Val()
}

func IsTrafficShift() bool {
	return manager.draining.Val()
}

func IsRejecting() bool {
	return manager.rejectingNewRequests.Val()
}

func StartTrafficShift() {
	manager.draining.Set(true)
	manager.rejectingNewRequests.Set(false)
}

func StartRejectNew() {
	manager.draining.Set(true)
	manager.rejectingNewRequests.Set(true)
}

func Resume() {
	manager.draining.Set(false)
	manager.rejectingNewRequests.Set(false)
}
