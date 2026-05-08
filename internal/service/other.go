package service

import (
	"context"

	"server_go/internal/model"
)

type IOther interface {
	GetResVersion(ctx context.Context, key string) (*model.ResVersionOutput, error)
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