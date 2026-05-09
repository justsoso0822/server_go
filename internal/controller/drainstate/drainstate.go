package drainstate

import "sync"

var manager = &stateManager{}

type stateManager struct {
	mu                   sync.RWMutex
	draining             bool
	rejectingNewRequests bool
	activeRequests       int64
}

func IncActiveRequests() {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.activeRequests++
}

func DecActiveRequests() {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.activeRequests--
}

func GetActiveRequests() int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.activeRequests
}

func IsTrafficShift() bool {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.draining
}

func IsRejecting() bool {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.rejectingNewRequests
}

func StartTrafficShift() {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.draining = true
	manager.rejectingNewRequests = false
}

func StartRejectNew() {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.draining = true
	manager.rejectingNewRequests = true
}

func Resume() {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.draining = false
	manager.rejectingNewRequests = false
}
