// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type (
	IOther interface {
		GetResVersion(ctx context.Context, key string) (g.Map, error)
	}
)

var (
	localOther IOther
)

func Other() IOther {
	if localOther == nil {
		panic("implement not found for interface IOther, forgot register?")
	}
	return localOther
}

func RegisterOther(i IOther) {
	localOther = i
}
