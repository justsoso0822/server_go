package control

import (
	"server_go/internal/runtime/drain"
	"server_go/internal/service"
)

type sControl struct{}

func init() {
	service.RegisterControl(&sControl{})
}

func (s *sControl) TrafficShift() (string, error) {
	drain.StartTrafficShift()
	return "traffic-shift", nil
}

func (s *sControl) RejectNew() (string, error) {
	drain.StartRejectNew()
	return "reject-new-requests", nil
}

func (s *sControl) ResumeTraffic() (string, error) {
	drain.Resume()
	return "resume-traffic", nil
}
