package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type ITask interface {
	InitTasks(ctx context.Context, uid int64) ([]g.Map, error)
}

var localTask ITask

func Task() ITask {
	if localTask == nil {
		panic("service ITask not registered")
	}
	return localTask
}

func RegisterTask(s ITask) {
	localTask = s
}