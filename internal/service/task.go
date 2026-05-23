// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ITask interface {
		InitTasks(ctx context.Context, uid int64) ([]map[string]interface{}, error)
	}
)

var (
	localTask ITask
)

func Task() ITask {
	if localTask == nil {
		panic("implement not found for interface ITask, forgot register?")
	}
	return localTask
}

func RegisterTask(i ITask) {
	localTask = i
}
