// Package dao provides database access helpers and query wrappers.
package dao

import (
	"context"

	"app/tools/autodb"

	"gorm.io/gorm"
)

func db(ctx context.Context) (*gorm.DB, error) {
	return autodb.DB(ctx)
}

// WithTx 在同一个数据库事务里执行 fn，把 *gorm.DB 直接透传给调用方。
//
//   - fn 返回 nil   → 提交
//   - fn 返回 error → 回滚，并把 error 原样返回
//   - fn panic      → 先回滚，再把 panic 抛出去
//
// 注意：fn 内必须用参数 tx 执行 SQL。若调用本包的其他方法（如 dao.GetUser），
// 它们走的是连接池，不在这条事务里，回滚也滚不到。
func WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(fn)
}
