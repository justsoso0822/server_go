package service

import (
	"context"

	"server_go/internal/model"
)

type IGame interface {
	Online(ctx context.Context, in *model.OnlineInput) error
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
