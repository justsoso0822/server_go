package dao

import (
	"context"
	"errors"
	"time"

	"server_go/dao/model"

	"gorm.io/gorm"
)

func GetUserBag(ctx context.Context, uid int64, chapter int) ([]model.UserBag, error) {
	b := q(ctx).UserBag
	rows, err := b.Where(b.UID.Eq(int32(uid)), b.Chapter.Eq(int32(chapter))).Find()
	return derefSlice(rows), err
}

func GetUserBagTp(ctx context.Context, uid int64, chapter int) ([]model.UserBagTp, error) {
	b := q(ctx).UserBagTp
	rows, err := b.Where(b.UID.Eq(int32(uid)), b.Chapter.Eq(int32(chapter))).Find()
	return derefSlice(rows), err
}

func GetUserTask(ctx context.Context, uid int64, minId, maxId int) (*model.UserTask, error) {
	t := q(ctx).UserTask
	row, err := t.
		Where(t.UID.Eq(int32(uid)), t.Taskid.Gte(int32(minId)), t.Taskid.Lte(int32(maxId))).
		Limit(1).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

func InsertUserTask(ctx context.Context, t *model.UserTask) error {
	return q(ctx).UserTask.Create(t)
}

func DeleteDoneUserTasks(ctx context.Context, uid int64, minId, maxId int) error {
	t := q(ctx).UserTask
	_, err := t.
		Where(t.UID.Eq(int32(uid)), t.Taskid.Gte(int32(minId)), t.Taskid.Lte(int32(maxId)), t.Done.Eq(1)).
		Delete()
	return err
}

// InsertUserOnline only inserts uid/day/tm_online and leaves tm_update to DB default.
// user_online.tm_update is NULL DEFAULT NULL; a full struct Create may write the
// zero time.Time value and exceed MySQL datetime's valid range.
func InsertUserOnline(ctx context.Context, o *model.UserOnline) error {
	uo := q(ctx).UserOnline
	return uo.Select(uo.UID, uo.Day, uo.TmOnline).Create(o)
}

// IncrUserOnlineTime atomically increments tm_online by delta and sets tm_update,
// returning the number of affected rows (0 means the row does not exist).
func IncrUserOnlineTime(ctx context.Context, uid int64, day time.Time, delta int32, tmUpdate time.Time) (int64, error) {
	o := q(ctx).UserOnline
	info, err := o.
		Where(o.UID.Eq(int32(uid)), o.Day.Eq(day)).
		Updates(map[string]interface{}{
			"tm_online": gorm.Expr("tm_online + ?", delta),
			"tm_update": tmUpdate,
		})
	if err != nil {
		return 0, err
	}
	return info.RowsAffected, nil
}

func GetPrfTaskMinMax(ctx context.Context, ser int) (minId, maxId int, err error) {
	t := q(ctx).PrfTask
	var res struct {
		Min int
		Max int
	}
	err = t.
		Select(t.ID.Min().As("min"), t.ID.Max().As("max")).
		Where(t.Ser.Eq(int32(ser))).
		Scan(&res)
	return res.Min, res.Max, err
}

func GetPrfTaskById(ctx context.Context, id int) (*model.PrfTask, error) {
	t := q(ctx).PrfTask
	row, err := t.Where(t.ID.Eq(int32(id))).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

func GetNextPrfTask(ctx context.Context, ser, afterId int) (int, error) {
	t := q(ctx).PrfTask
	row, err := t.Select(t.ID).
		Where(t.Ser.Eq(int32(ser)), t.ID.Gt(int32(afterId))).
		Order(t.ID).Limit(1).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int(row.ID), nil
}

func GetLoopStartPrfTask(ctx context.Context, ser int) (int, error) {
	t := q(ctx).PrfTask
	row, err := t.Select(t.ID).
		Where(t.Ser.Eq(int32(ser)), t.StartLoop.Eq(1)).
		Order(t.ID).Limit(1).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int(row.ID), nil
}
