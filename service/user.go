package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"server_go/dao"
	"server_go/dao/model"
	"server_go/state"
	"server_go/tools/autodb"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func UserLogin(ctx context.Context, uid int64, loginKey, openid, platform, version string) (map[string]any, error) {
	if openid == "" {
		return nil, fmt.Errorf("参数错误: openid 必填")
	}

	lockKey := fmt.Sprintf("user_login:%d", uid)
	token, err := state.Lock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	defer func() { _ = state.Unlock(ctx, lockKey, token) }()

	out := map[string]any{"uid": uid}

	user, err := dao.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if user.Platform != platform || user.Openid != openid {
			return nil, fmt.Errorf("账号信息不匹配")
		}
		out["newbie"] = 0
		out["user"] = user
	} else {
		out["newbie"] = 1
		nowDay := int32(startOfDay(time.Now()).Unix())
		err = autodb.DB(ctx).Transaction(func(tx *gorm.DB) error {
			// GORM Transaction 会在回调返回 nil 时提交，返回 error/panic 时回滚。
			if e := tx.Create(&model.User{UID: uid, Platform: platform, Openid: openid}).Error; e != nil {
				return e
			}
			return tx.Create(&model.UserRes{
				UID: uid, Gold: 200, Diamond: 100, Star: 0,
				Tili: 100, TiliTime: 0, Exp: 0, Level: 1, DayTime: nowDay,
			}).Error
		})
		if err != nil {
			return nil, err
		}
		created, err := dao.GetUser(ctx, uid)
		if err != nil {
			return nil, err
		}
		out["user"] = created
	}

	bgCtx := autodb.BackgroundWithChannel(ctx)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logAsyncPanic(bgCtx, "async login log panic", r,
					zap.Int64("uid", uid),
					zap.String("platform", platform),
				)
			}
		}()
		tCtx, cancel := context.WithTimeout(bgCtx, asyncLogTimeout)
		defer cancel()
		if err := dao.InsertLogLogin(tCtx, &model.LogLogin{
			UID:       uid,
			Platform:  platform,
			RequestID: autodb.GetRequestID(bgCtx),
		}); err != nil {
			logAsyncError(bgCtx, "async login log failed", err,
				zap.Int64("uid", uid),
				zap.String("platform", platform),
			)
		}
	}()

	ver, _ := strconv.ParseInt(version, 10, 32)
	if err = dao.SaveUserLoginkey(ctx, &model.UserLoginkey{
		UID: uid, Key: loginKey, Ver: int32(ver), Time: int32(time.Now().Unix()),
	}); err != nil {
		return nil, err
	}

	rc := autodb.Redis(ctx)
	cacheKey := state.BuildKey(ctx, "login_key", "uid", strconv.FormatInt(uid, 10))
	if err = rc.SetEx(ctx, cacheKey, loginKey, 7200*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("cache login_key: %w", err)
	}

	datas, err := dao.GetUserDatas(ctx, uid)
	if err != nil {
		return nil, err
	}
	out["datas"] = datas

	gm, err := dao.IsGm(ctx, uid)
	if err != nil {
		return nil, err
	}
	if gm {
		out["gm"] = 1
	} else {
		out["gm"] = 0
	}

	items, err := dao.GetUserItems(ctx, uid)
	if err != nil {
		return nil, err
	}
	out["items"] = items

	res, err := GetUserRes(ctx, uid)
	if err != nil {
		return nil, err
	}
	out["res"] = res

	configs, err := dao.GetAllMemConfig(ctx)
	if err != nil {
		return nil, err
	}
	out["config"] = configs

	return out, nil
}

func GetUser(ctx context.Context, uid int64) (*model.User, error) {
	return dao.GetUser(ctx, uid)
}

func GetUserRes(ctx context.Context, uid int64) (*model.UserRes, error) {
	res, err := dao.GetUserRes(ctx, uid)
	if err != nil || res == nil {
		return res, err
	}
	nowDay := int32(startOfDay(time.Now()).Unix())
	if res.DayTime != nowDay {
		_ = dao.UpdateUserResDayConf(ctx, uid, "", nowDay)
		res.DayConf = ""
		res.DayTime = nowDay
	}
	return res, nil
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
