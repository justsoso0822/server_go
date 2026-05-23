// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IControl interface {
		TrafficShift() (string, error)
		RejectNew() (string, error)
		ResumeTraffic() (string, error)
	}
)

var (
	localControl IControl
)

func Control() IControl {
	if localControl == nil {
		panic("implement not found for interface IControl, forgot register?")
	}
	return localControl
}

func RegisterControl(i IControl) {
	localControl = i
}
