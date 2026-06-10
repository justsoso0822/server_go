package dao

import (
	"context"
	"errors"
	"fmt"
	"math"

	"server_go/dao/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUser(ctx context.Context, uid int64) (*model.User, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var row model.User
	err = db.Where("uid = ?", uid).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &row, err
}

func InsertUser(ctx context.Context, u *model.User) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Create(u).Error
}

func GetUserRes(ctx context.Context, uid int64) (*model.UserRes, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var row model.UserRes
	err = db.Where("uid = ?", uid).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &row, err
}

func InsertUserRes(ctx context.Context, r *model.UserRes) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Create(r).Error
}

type UserResField string

const (
	UserResFieldDiamond UserResField = "diamond"
	UserResFieldGold    UserResField = "gold"
	UserResFieldTili    UserResField = "tili"
	UserResFieldExp     UserResField = "exp"
	UserResFieldStar    UserResField = "star"
)

func UpdateUserResField(ctx context.Context, uid int64, resField UserResField, value int64) error {
	colName := string(resField)
	switch resField {
	case UserResFieldDiamond, UserResFieldGold, UserResFieldTili, UserResFieldExp, UserResFieldStar:
	default:
		return errors.New("invalid user resource field")
	}
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Model(&model.UserRes{}).Where("uid = ?", uid).Update(colName, value).Error
}

// 原子地增减用户资源，delta 为正表示增加、负表示减少。
// 即使上层锁失效，数据库层面也会通过 WHERE 守卫确保余额不为负且不超 int32 上限，
// 并发写入是安全的：MySQL UPDATE 会对匹配的行加行锁串行执行。
func IncrUserResField(ctx context.Context, uid int64, field UserResField, delta int64) error {
	colName := string(field)
	db, err := db(ctx)
	if err != nil {
		return err
	}

	result := db.
		Model(&model.UserRes{}).
		Where("uid = ?", uid).
		Where(colName+" + ? >= 0", delta).
		Where(colName+" + ? <= ?", delta, math.MaxInt32).
		// gorm.Expr 会生成 column = column + ?，避免先读后写导致的并发覆盖。
		Update(colName, gorm.Expr(colName+" + ?", delta))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("资源变更失败")
	}
	return nil
}

func UpdateUserResDayConf(ctx context.Context, uid int64, dayConf string, dayTime int32) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Model(&model.UserRes{}).
		Where("uid = ?", uid).
		Updates(map[string]interface{}{
			"day_conf": dayConf,
			"day_time": dayTime,
		}).Error
}

func GetUserLoginkey(ctx context.Context, uid int64, key string) (*model.UserLoginkey, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var row model.UserLoginkey
	err = db.Where("uid = ? AND `key` = ?", uid, key).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &row, err
}

func SaveUserLoginkey(ctx context.Context, k *model.UserLoginkey) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	// OnConflict{UpdateAll:true} 是 GORM 的 upsert 写法；
	// MySQL 下会生成 ON DUPLICATE KEY UPDATE，依赖表上的唯一键/主键判断冲突。
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(k).Error
}

func GetUserDatas(ctx context.Context, uid int64) ([]model.UserData, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []model.UserData
	err = db.Where("uid = ?", uid).Find(&rows).Error
	return rows, err
}

func GetUserDataValue(ctx context.Context, uid int64, key string) (string, error) {
	db, err := db(ctx)
	if err != nil {
		return "", err
	}
	var row model.UserData
	err = db.Where("uid = ? AND `key` = ?", uid, key).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

func InsertUserData(ctx context.Context, d *model.UserData) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Create(d).Error
}

func GetUserItems(ctx context.Context, uid int64) ([]model.UserItem, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []model.UserItem
	err = db.Where("uid = ?", uid).Find(&rows).Error
	return rows, err
}
