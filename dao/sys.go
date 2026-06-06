package dao

import (
	"context"
	"errors"

	"server_go/dao/model"

	"gorm.io/gorm"
)

func GetAllMemConfig(ctx context.Context) ([]model.MemConfig, error) {
	var rows []model.MemConfig
	err := db(ctx).Find(&rows).Error
	return rows, err
}

func GetMemConfigValue(ctx context.Context, id int) (string, error) {
	var row model.MemConfig
	err := db(ctx).Where("id = ?", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

func IsGm(ctx context.Context, uid int64) (bool, error) {
	var count int64
	err := db(ctx).Model(&model.SysGm{}).Where("uid = ?", uid).Count(&count).Error
	return count > 0, err
}
