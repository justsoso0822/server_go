package service

import (
	"context"
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
	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	nowStr := now.Format("2006-01-02 15:04:05")

	row, err := dao.GetUserOnline(ctx, uid, day.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	if row != nil {
		seconds += int64(row.TmOnline)
		return dao.UpdateUserOnline(ctx, uid, day.Format("2006-01-02 15:04:05"), seconds, nowStr)
	}
	return dao.InsertUserOnline(ctx, &model.UserOnline{
		UID: int32(uid), Day: day, TmOnline: int32(seconds),
	})
}
