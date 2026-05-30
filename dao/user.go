package dao

import (
	"context"
	"errors"

	"server_gin/dao/model"

	"gorm.io/gorm"
)

func GetUser(ctx context.Context, uid int64) (*model.User, error) {
	var u model.User
	err := db(ctx).Where("uid = ?", uid).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func InsertUser(ctx context.Context, u *model.User) error {
	return db(ctx).Create(u).Error
}

func GetUserRes(ctx context.Context, uid int64) (*model.UserRes, error) {
	var r model.UserRes
	err := db(ctx).Where("uid = ?", uid).First(&r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &r, err
}

func InsertUserRes(ctx context.Context, r *model.UserRes) error {
	return db(ctx).Create(r).Error
}

type UserResField string

const (
	UserResFieldDiamond UserResField = "diamond"
	UserResFieldGold    UserResField = "gold"
	UserResFieldTili    UserResField = "tili"
	UserResFieldExp     UserResField = "exp"
	UserResFieldStar    UserResField = "star"
)

var userResColumns = map[UserResField]string{
	UserResFieldDiamond: "diamond",
	UserResFieldGold:    "gold",
	UserResFieldTili:    "tili",
	UserResFieldExp:     "exp",
	UserResFieldStar:    "star",
}

func UpdateUserResField(ctx context.Context, uid int64, field UserResField, value int64) error {
	column, ok := userResColumns[field]
	if !ok {
		return errors.New("invalid user resource field")
	}
	return db(ctx).Model(&model.UserRes{}).Where("uid = ?", uid).Update(column, value).Error
}

func UpdateUserResDayConf(ctx context.Context, uid int64, dayConf string, dayTime int32) error {
	return db(ctx).Model(&model.UserRes{}).Where("uid = ?", uid).
		Updates(map[string]any{"day_conf": dayConf, "day_time": dayTime}).Error
}

func GetUserLoginkey(ctx context.Context, uid int64, key string) (*model.UserLoginkey, error) {
	var k model.UserLoginkey
	err := db(ctx).Where("uid = ? AND `key` = ?", uid, key).First(&k).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &k, err
}

func SaveUserLoginkey(ctx context.Context, k *model.UserLoginkey) error {
	return db(ctx).Save(k).Error
}

func GetUserDatas(ctx context.Context, uid int64) ([]model.UserData, error) {
	var rows []model.UserData
	err := db(ctx).Where("uid = ?", uid).Find(&rows).Error
	return rows, err
}

func GetUserDataValue(ctx context.Context, uid int64, key string) (string, error) {
	var d model.UserData
	err := db(ctx).Where("uid = ? AND `key` = ?", uid, key).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return d.Value, err
}

func InsertUserData(ctx context.Context, d *model.UserData) error {
	return db(ctx).Create(d).Error
}

func GetUserItems(ctx context.Context, uid int64) ([]model.UserItem, error) {
	var rows []model.UserItem
	err := db(ctx).Where("uid = ?", uid).Find(&rows).Error
	return rows, err
}
