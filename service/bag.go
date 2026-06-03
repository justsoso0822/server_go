package service

import (
	"context"
	"math"
	"time"

	"server_go/dao"
	"server_go/dao/model"
)

func GetUserBag(ctx context.Context, uid int64, chapter int) (map[string]any, error) {
	rows, err := dao.GetUserBag(ctx, uid, chapter)
	if err != nil {
		return nil, err
	}
	return map[string]any{"uid": uid, "chapter": chapter, "bag": rows}, nil
}

func GetUserBagTp(ctx context.Context, uid int64, chapter int) (map[string]any, error) {
	rows, err := dao.GetUserBagTp(ctx, uid, chapter)
	if err != nil {
		return nil, err
	}
	return map[string]any{"uid": uid, "chapter": chapter, "bag": rows}, nil
}

func GameOnline(ctx context.Context, uid int64, seconds int64) error {
	if seconds > math.MaxInt32 {
		seconds = math.MaxInt32
	}
	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// Atomic increment: UPDATE ... SET tm_online = tm_online + ? WHERE uid = ? AND day = ?
	affected, err := dao.IncrUserOnlineTime(ctx, uid, day, int32(seconds), now)
	if err != nil {
		return err
	}
	if affected > 0 {
		return nil
	}
	// Row doesn't exist yet, insert it.
	return dao.InsertUserOnline(ctx, &model.UserOnline{
		UID: uid, Day: day, TmOnline: int32(seconds),
	})
}
