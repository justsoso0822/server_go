package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type IOther interface {
	GetResVersion(ctx context.Context, key string) (g.Map, error)
}

var localOther IOther

func Other() IOther {
	if localOther == nil {
		panic("service IOther not registered")
	}
	return localOther
}

func RegisterOther(s IOther) {
	localOther = s
}