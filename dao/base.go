// Package dao provides database access helpers and query wrappers.
package dao

import (
	"context"

	"server_go/tools/autodb"

	"gorm.io/gorm"
)

func db(ctx context.Context) (*gorm.DB, error) {
	return autodb.DB(ctx)
}
