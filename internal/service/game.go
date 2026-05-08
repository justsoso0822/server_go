package service

import "context"

type IGame interface {
	Online(ctx context.Context, uid, seconds int64) (map[string]interface{}, error)
	ServerTime() map[string]interface{}
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