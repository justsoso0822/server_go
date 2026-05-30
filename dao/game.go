package dao

import (
	"context"
	"errors"

	"server_gin/dao/model"

	"gorm.io/gorm"
)

func GetUserBag(ctx context.Context, uid int64, chapter int) ([]model.UserBag, error) {
	var rows []model.UserBag
	err := db(ctx).Where("uid = ? AND chapter = ?", uid, chapter).Find(&rows).Error
	return rows, err
}

func GetUserBagTp(ctx context.Context, uid int64, chapter int) ([]model.UserBagTp, error) {
	var rows []model.UserBagTp
	err := db(ctx).Where("uid = ? AND chapter = ?", uid, chapter).Find(&rows).Error
	return rows, err
}

func GetUserTask(ctx context.Context, uid int64, minId, maxId int) (*model.UserTask, error) {
	var t model.UserTask
	err := db(ctx).Where("uid = ? AND taskid >= ? AND taskid <= ?", uid, minId, maxId).
		Limit(1).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func InsertUserTask(ctx context.Context, t *model.UserTask) error {
	return db(ctx).Create(t).Error
}

func DeleteDoneUserTasks(ctx context.Context, uid int64, minId, maxId int) error {
	return db(ctx).Where("uid = ? AND taskid >= ? AND taskid <= ? AND done = 1", uid, minId, maxId).
		Delete(&model.UserTask{}).Error
}

func GetUserOnline(ctx context.Context, uid int64, day string) (*model.UserOnline, error) {
	var o model.UserOnline
	err := db(ctx).Where("uid = ? AND day = ?", uid, day).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &o, err
}

func InsertUserOnline(ctx context.Context, o *model.UserOnline) error {
	return db(ctx).Table("user_online").Create(map[string]interface{}{
		"uid": o.UID, "day": o.Day, "tm_online": o.TmOnline,
	}).Error
}

func UpdateUserOnline(ctx context.Context, uid int64, day string, tmOnline int64, tmUpdate string) error {
	return db(ctx).Model(&model.UserOnline{}).
		Where("uid = ? AND day = ?", uid, day).
		Updates(map[string]interface{}{"tm_online": tmOnline, "tm_update": tmUpdate}).Error
}

func GetPrfTaskMinMax(ctx context.Context, ser int) (minId, maxId int, err error) {
	var res struct {
		Min int
		Max int
	}
	err = db(ctx).Model(&model.PrfTask{}).Select("MIN(id) as min, MAX(id) as max").
		Where("ser = ?", ser).Scan(&res).Error
	return res.Min, res.Max, err
}

func GetPrfTaskById(ctx context.Context, id int) (*model.PrfTask, error) {
	var t model.PrfTask
	err := db(ctx).Where("id = ?", id).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func GetNextPrfTask(ctx context.Context, ser, afterId int) (int, error) {
	var t model.PrfTask
	err := db(ctx).Select("id").Where("ser = ? AND id > ?", ser, afterId).
		Order("id ASC").Limit(1).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return int(t.ID), err
}

func GetLoopStartPrfTask(ctx context.Context, ser int) (int, error) {
	var t model.PrfTask
	err := db(ctx).Select("id").Where("ser = ? AND start_loop = 1", ser).
		Order("id ASC").Limit(1).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return int(t.ID), err
}
