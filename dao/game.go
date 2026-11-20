package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"server_go/dao/model"
	"server_go/tools/autodb"

	"gorm.io/gorm"
)

func GetUserBag(ctx context.Context, uid int64, chapter int) ([]model.UserBag, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []model.UserBag
	err = db.Where("uid = ? AND chapter = ?", uid, chapter).Find(&rows).Error
	return rows, err
}

func GetUserBagTp(ctx context.Context, uid int64, chapter int) ([]model.UserBagTp, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var rows []model.UserBagTp
	err = db.Where("uid = ? AND chapter = ?", uid, chapter).Find(&rows).Error
	return rows, err
}

func GetUserTask(ctx context.Context, uid int64, minID, maxID int) (*model.UserTask, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var row model.UserTask
	err = db.
		Where("uid = ? AND taskid >= ? AND taskid <= ?", uid, minID, maxID).
		Limit(1).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &row, err
}

func InsertUserTask(ctx context.Context, t *model.UserTask) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Create(t).Error
}

func DeleteDoneUserTasks(ctx context.Context, uid int64, minID, maxID int) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.
		Where("uid = ? AND taskid >= ? AND taskid <= ? AND done = ?", uid, minID, maxID, 1).
		Delete(&model.UserTask{}).Error
}

// InsertUserOnline only inserts uid/day/tm_online and leaves tm_update to DB default.
// user_online.tm_update is NULL DEFAULT NULL; a full struct Create may write the
// zero time.Time value and exceed MySQL datetime's valid range.
func InsertUserOnline(ctx context.Context, o *model.UserOnline) error {
	db, err := db(ctx)
	if err != nil {
		return err
	}
	return db.Select("uid", "day", "tm_online").Create(o).Error
}

// IncrUserOnlineTime atomically increments tm_online by delta and sets tm_update,
// returning the number of affected rows (0 means the row does not exist).
func IncrUserOnlineTime(ctx context.Context, uid int64, day time.Time, delta int32, tmUpdate time.Time) (int64, error) {
	db, err := db(ctx)
	if err != nil {
		return 0, err
	}
	result := db.Model(&model.UserOnline{}).
		Where("uid = ? AND day = ?", uid, day).
		Updates(map[string]any{
			"tm_online": gorm.Expr("tm_online + ?", delta),
			"tm_update": tmUpdate,
		})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func GetPrfTaskMinMax(ctx context.Context, ser int) (minID, maxID int, err error) {
	type prfTaskMinMax struct {
		Min int
		Max int
	}
	res, err := autodb.Cache(ctx, fmt.Sprintf("prf_task_min_max:%d", ser), time.Minute, func() (prfTaskMinMax, error) {
		db, err := db(ctx)
		if err != nil {
			return prfTaskMinMax{}, err
		}
		var r prfTaskMinMax
		err = db.Model(&model.PrfTask{}).
			Select("COALESCE(MIN(id), 0) AS min, COALESCE(MAX(id), 0) AS max").
			Where("ser = ?", ser).
			Scan(&r).Error
		return r, err
	})
	return res.Min, res.Max, err
}

func GetPrfTaskByID(ctx context.Context, id int) (*model.PrfTask, error) {
	db, err := db(ctx)
	if err != nil {
		return nil, err
	}
	var row model.PrfTask
	err = db.Where("id = ?", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &row, err
}

func GetNextPrfTask(ctx context.Context, ser, afterID int) (int, error) {
	db, err := db(ctx)
	if err != nil {
		return 0, err
	}
	var row model.PrfTask
	err = db.Select("id").
		Where("ser = ? AND id > ?", ser, afterID).
		Order("id").Limit(1).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int(row.ID), nil
}

func GetLoopStartPrfTask(ctx context.Context, ser int) (int, error) {
	db, err := db(ctx)
	if err != nil {
		return 0, err
	}
	var row model.PrfTask
	err = db.Select("id").
		Where("ser = ? AND start_loop = ?", ser, 1).
		Order("id").Limit(1).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int(row.ID), nil
}
