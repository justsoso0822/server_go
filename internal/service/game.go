package service

import (
	"context"
)

type IGame interface {
	Online(ctx context.Context, uid int64, seconds int64) error
}

var localGame IGame

func Game() IGame {
	if localGame == nil {
		panic("service IGame not registered")
	}
	return localGame
}

func RegisterGame(s IGame) {
	localGame = s
}
