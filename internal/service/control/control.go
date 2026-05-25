package control

import "server_go/internal/runtime/drain"

func TrafficShift() (string, error) {
	drain.StartTrafficShift()
	return "traffic-shift", nil
}

func RejectNew() (string, error) {
	drain.StartRejectNew()
	return "reject-new-requests", nil
}

func ResumeTraffic() (string, error) {
	drain.Resume()
	return "resume-traffic", nil
}
