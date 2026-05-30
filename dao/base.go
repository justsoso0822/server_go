package dao

import (
	"context"

	"server_gin/autodb"

	"gorm.io/gorm"
)

func db(ctx context.Context) *gorm.DB {
	return autodb.DB(ctx)
}
