package dao

import (
	"context"
	"errors"

	"server_go/dao/model"

	"gorm.io/gorm"
)

func GetUserBag(ctx context.Context, uid int64, chapter int) ([]model.UserBag, error) {
	b := q(ctx).UserBag
	rows, err := b.WithContext(ctx).Where(b.UID.Eq(int32(uid)), b.Chapter.Eq(int32(chapter))).Find()
	return derefSlice(rows), err
}

func GetUserBagTp(ctx context.Context, uid int64, chapter int) ([]model.UserBagTp, error) {
	b := q(ctx).UserBagTp
	rows, err := b.WithContext(ctx).Where(b.UID.Eq(int32(uid)), b.Chapter.Eq(int32(chapter))).Find()
	return derefSlice(rows), err
}

func GetUserTask(ctx context.Context, uid int64, minId, maxId int) (*model.UserTask, error) {
	t := q(ctx).UserTask
	row, err := t.WithContext(ctx).
		Where(t.UID.Eq(int32(uid)), t.Taskid.Gte(int32(minId)), t.Taskid.Lte(int32(maxId))).
		Limit(1).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

func InsertUserTask(ctx context.Context, t *model.UserTask) error {
	return q(ctx).UserTask.WithContext(ctx).Create(t)
}

func DeleteDoneUserTasks(ctx context.Context, uid int64, minId, maxId int) error {
	t := q(ctx).UserTask
	_, err := t.WithContext(ctx).
		Where(t.UID.Eq(int32(uid)), t.Taskid.Gte(int32(minId)), t.Taskid.Lte(int32(maxId)), t.Done.Eq(1)).
		Delete()
	return err
}

func GetUserOnline(ctx context.Context, uid int64, day string) (*model.UserOnline, error) {
	o := q(ctx).UserOnline
	row, err := o.WithContext(ctx).Where(o.UID.Eq(int32(uid)), o.Day.Eq(parseDateTime(day))).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

// InsertUserOnline 保留原始 map + Table 写法：user_online.tm_update 为
// NULL DEFAULT NULL 且 model 无 default 标签，若用 struct Create，未赋值的
// time.Time 零值（0001-01-01）会触发 MySQL datetime 越界。这里只插入 3 列，
// 让 tm_update 走 DB 默认值（NULL）。
func InsertUserOnline(ctx context.Context, o *model.UserOnline) error {
	return db(ctx).Table("user_online").Create(map[string]interface{}{
		"uid": o.UID, "day": o.Day, "tm_online": o.TmOnline,
	}).Error
}

func UpdateUserOnline(ctx context.Context, uid int64, day string, tmOnline int64, tmUpdate string) error {
	o := q(ctx).UserOnline
	_, err := o.WithContext(ctx).
		Where(o.UID.Eq(int32(uid)), o.Day.Eq(parseDateTime(day))).
		Updates(map[string]interface{}{"tm_online": tmOnline, "tm_update": tmUpdate})
	return err
}

func GetPrfTaskMinMax(ctx context.Context, ser int) (minId, maxId int, err error) {
	t := q(ctx).PrfTask
	var res struct {
		Min int
		Max int
	}
	err = t.WithContext(ctx).
		Select(t.ID.Min().As("min"), t.ID.Max().As("max")).
		Where(t.Ser.Eq(int32(ser))).
		Scan(&res)
	return res.Min, res.Max, err
}

func GetPrfTaskById(ctx context.Context, id int) (*model.PrfTask, error) {
	t := q(ctx).PrfTask
	row, err := t.WithContext(ctx).Where(t.ID.Eq(int32(id))).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return row, err
}

func GetNextPrfTask(ctx context.Context, ser, afterId int) (int, error) {
	t := q(ctx).PrfTask
	row, err := t.WithContext(ctx).Select(t.ID).
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
	row, err := t.WithContext(ctx).Select(t.ID).
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

