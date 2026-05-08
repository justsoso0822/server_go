package service

import (
	"context"
)

type ITask interface {
	InitTasks(ctx context.Context, uid int64) ([]map[string]interface{}, error)
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