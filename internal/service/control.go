package service

type IControl interface {
	TrafficShift() (string, error)
	RejectNew() (string, error)
	ResumeTraffic() (string, error)
}

var localControl IControl

func Control() IControl {
	if localControl == nil {
		panic("service IControl not registered")
	}
	return localControl
}

func RegisterControl(s IControl) {
	localControl = s
}
