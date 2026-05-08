package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type IBag interface {
	GetUserBag(ctx context.Context, uid int64, chapter int) (g.Map, error)
	GetUserBagTp(ctx context.Context, uid int64, chapter int) (g.Map, error)
}

var localBag IBag

func Bag() IBag {
	if localBag == nil {
		panic("service IBag not registered")
	}
	return localBag
}

func RegisterBag(s IBag) {
	localBag = s
}