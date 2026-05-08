package service

import (
	"context"

	"server_go/internal/model"
)

type IBag interface {
	GetUserBag(ctx context.Context, in *model.BagInput) (*model.BagOutput, error)
	GetUserBagTp(ctx context.Context, in *model.BagInput) (*model.BagOutput, error)
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