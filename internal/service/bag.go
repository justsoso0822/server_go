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
	IBag interface {
		GetUserBag(ctx context.Context, uid int64, chapter int) (g.Map, error)
		GetUserBagTp(ctx context.Context, uid int64, chapter int) (g.Map, error)
	}
)

var (
	localBag IBag
)

func Bag() IBag {
	if localBag == nil {
		panic("implement not found for interface IBag, forgot register?")
	}
	return localBag
}

func RegisterBag(i IBag) {
	localBag = i
}
