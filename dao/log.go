package dao

import (
	"context"

	"server_go/dao/model"
)

func InsertLogLogin(ctx context.Context, l *model.LogLogin) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Create(l).Error
}

func InsertLogTrace(ctx context.Context, l *model.LogTrace) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Create(l).Error
}

func InsertLogMsg(ctx context.Context, l *model.LogMsg) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Create(l).Error
}
