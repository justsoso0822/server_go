package dao

import (
	"context"
	"errors"

	"server_go/dao/model"

	"gorm.io/gorm"
)

func GetAllMemConfig(ctx context.Context) ([]model.MemConfig, error) {
	rows, err := q(ctx).MemConfig.Find()
	return derefSlice(rows), err
}

func GetMemConfigValue(ctx context.Context, id int) (string, error) {
	c := q(ctx).MemConfig
	row, err := c.Where(c.ID.Eq(int32(id))).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

func IsGm(ctx context.Context, uid int64) (bool, error) {
	g := q(ctx).SysGm
	count, err := g.Where(g.UID.Eq(int32(uid))).Count()
	return count > 0, err
}
