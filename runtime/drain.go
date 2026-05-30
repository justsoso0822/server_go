package runtime

import "sync/atomic"

var (
	draining             atomic.Bool
	rejectingNewRequests atomic.Bool
	activeRequests       atomic.Int64
)

func IncActiveRequests() { activeRequests.Add(1) }

func DecActiveRequests() {
	if v := activeRequests.Add(-1); v < 0 {
		activeRequests.Store(0)
	}
}

func GetActiveRequests() int64 { return activeRequests.Load() }
func IsTrafficShift() bool     { return draining.Load() }
func IsRejecting() bool        { return rejectingNewRequests.Load() }

func StartTrafficShift() { draining.Store(true); rejectingNewRequests.Store(false) }
func StartRejectNew()    { draining.Store(true); rejectingNewRequests.Store(true) }
func Resume()            { draining.Store(false); rejectingNewRequests.Store(false) }
