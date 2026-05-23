package control

import (
	"server_go/internal/logic/drainstate"
	"server_go/internal/service"
)

type sControl struct{}

func init() {
	service.RegisterControl(&sControl{})
}

func (s *sControl) TrafficShift() (string, error) {
	drainstate.StartTrafficShift()
	return "traffic-shift", nil
}

func (s *sControl) RejectNew() (string, error) {
	drainstate.StartRejectNew()
	return "reject-new-requests", nil
}

func (s *sControl) ResumeTraffic() (string, error) {
	drainstate.Resume()
	return "resume-traffic", nil
}
