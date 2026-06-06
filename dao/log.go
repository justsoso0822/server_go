package dao

import (
	"context"

	"server_go/dao/model"
)

func InsertLogLogin(ctx context.Context, l *model.LogLogin) error {
	return db(ctx).Create(l).Error
}

func InsertLogTrace(ctx context.Context, l *model.LogTrace) error {
	return db(ctx).Create(l).Error
}

func InsertLogMsg(ctx context.Context, l *model.LogMsg) error {
	return db(ctx).Create(l).Error
}
