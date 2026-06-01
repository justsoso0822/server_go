package dao

import (
	"context"
	"errors"

	"server_go/dao/model"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func GetUser(ctx context.Context, uid int64) (*model.User, error) {
	u := q(ctx).User
	row, err := u.Where(u.UID.Eq(int32(uid))).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

func InsertUser(ctx context.Context, u *model.User) error {
	return q(ctx).User.Create(u)
}

func GetUserRes(ctx context.Context, uid int64) (*model.UserRes, error) {
	r := q(ctx).UserRes
	row, err := r.Where(r.UID.Eq(int32(uid))).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

func InsertUserRes(ctx context.Context, r *model.UserRes) error {
	return q(ctx).UserRes.Create(r)
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
	r := q(ctx).UserRes
	column, ok := map[UserResField]field.Expr{
		UserResFieldDiamond: r.Diamond,
		UserResFieldGold:    r.Gold,
		UserResFieldTili:    r.Tili,
		UserResFieldExp:     r.Exp,
		UserResFieldStar:    r.Star,
	}[resField]
	if !ok {
		return errors.New("invalid user resource field")
	}
	_, err := r.Where(r.UID.Eq(int32(uid))).Update(column, value)
	return err
}

func UpdateUserResDayConf(ctx context.Context, uid int64, dayConf string, dayTime int32) error {
	r := q(ctx).UserRes
	_, err := r.Where(r.UID.Eq(int32(uid))).
		UpdateSimple(r.DayConf.Value(dayConf), r.DayTime.Value(dayTime))
	return err
}

func GetUserLoginkey(ctx context.Context, uid int64, key string) (*model.UserLoginkey, error) {
	k := q(ctx).UserLoginkey
	row, err := k.Where(k.UID.Eq(int32(uid)), k.Key.Eq(key)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

func SaveUserLoginkey(ctx context.Context, k *model.UserLoginkey) error {
	return q(ctx).UserLoginkey.Save(k)
}

func GetUserDatas(ctx context.Context, uid int64) ([]model.UserData, error) {
	d := q(ctx).UserData
	rows, err := d.Where(d.UID.Eq(int32(uid))).Find()
	return derefSlice(rows), err
}

func GetUserDataValue(ctx context.Context, uid int64, key string) (string, error) {
	d := q(ctx).UserData
	row, err := d.Where(d.UID.Eq(int32(uid)), d.Key.Eq(key)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.Value, nil
}

func InsertUserData(ctx context.Context, d *model.UserData) error {
	return q(ctx).UserData.Create(d)
}

func GetUserItems(ctx context.Context, uid int64) ([]model.UserItem, error) {
	i := q(ctx).UserItem
	rows, err := i.Where(i.UID.Eq(int32(uid))).Find()
	return derefSlice(rows), err
}
