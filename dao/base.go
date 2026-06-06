// Package dao provides database access helpers and query wrappers.
package dao

import (
	"context"
	"time"

	"server_go/dao/query"
	"server_go/tools/autodb"

	"gorm.io/gorm"
)

func db(ctx context.Context) *gorm.DB {
	return autodb.DB(ctx)
}

// q 基于当前请求渠道的 *gorm.DB 构造 gorm gen 查询入口。
// 多渠道架构下每个 channel 持有独立连接，因此不能使用 query 包的全局
// SetDefault/Q，必须每次按 ctx 取连接构造。
func q(ctx context.Context) *query.Query {
	return query.Use(autodb.DB(ctx))
}

// derefSlice 将 gorm gen 返回的 []*T 转为调用方期望的 []T。
func derefSlice[T any](in []*T) []T {
	if len(in) == 0 {
		return nil
	}
	out := make([]T, 0, len(in))
	for _, p := range in {
		if p != nil {
			out = append(out, *p)
		}
	}
	return out
}

// parseDateTime 将 "2006-01-02 15:04:05" 文本解析为 Asia/Shanghai 时区 time.Time，
// 供 gorm gen 的强类型 datetime 字段条件使用。解析失败返回零值。
func parseDateTime(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, shanghai)
	if err != nil {
		return time.Time{}
	}
	return t
}

var shanghai *time.Location

func init() {
	shanghai, _ = time.LoadLocation("Asia/Shanghai")
}
