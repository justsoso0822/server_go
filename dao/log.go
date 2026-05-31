package dao

import (
	"context"

	"server_go/dao/model"
)

func InsertLogLogin(ctx context.Context, l *model.LogLogin) error {
	return q(ctx).LogLogin.WithContext(ctx).Create(l)
}

func InsertLogTrace(ctx context.Context, l *model.LogTrace) error {
	return q(ctx).LogTrace.WithContext(ctx).Create(l)
}

func InsertLogMsg(ctx context.Context, l *model.LogMsg) error {
	return q(ctx).LogMsg.WithContext(ctx).Create(l)
}
