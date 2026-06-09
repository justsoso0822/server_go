// 排水状态管理，配合部署脚本实现平滑切流和回滚。
package state

import "sync/atomic"

var (
	draining             atomic.Bool
	rejectingNewRequests atomic.Bool
	activeRequests       atomic.Int64
)

// 业务请求进入 DrainGuard 后加一。atomic 避免每个请求都抢 mutex。
func IncActiveRequests() { activeRequests.Add(1) }

// 业务请求结束后减一；防御性归零是为了避免异常调用顺序导致计数长期为负。
func DecActiveRequests() {
	if v := activeRequests.Add(-1); v < 0 {
		activeRequests.Store(0)
	}
}

// 给健康检查详情和部署脚本读取当前在途请求数。
func GetActiveRequests() int64 { return activeRequests.Load() }

// 只表示“正在切流/已摘流”，不等于已经拒绝新请求。
func IsTrafficShift() bool { return draining.Load() }

// 表示业务入口已经开始拒绝新请求。
func IsRejecting() bool { return rejectingNewRequests.Load() }

// 第一阶段：通知负载均衡摘流，但暂时不拒绝请求，给网关切换留出确认窗口。
func StartTrafficShift() { draining.Store(true); rejectingNewRequests.Store(false) }

// 第二阶段：确认流量已切走后，拒绝落到旧实例的新请求，只等待存量请求结束。
func StartRejectNew() { draining.Store(true); rejectingNewRequests.Store(true) }

// 手动恢复或回滚时清空排水状态，让实例重新对外提供服务。
func Resume() { draining.Store(false); rejectingNewRequests.Store(false) }
