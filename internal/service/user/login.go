package user

import (
	"context"
	"fmt"
	"strconv"

	"server_go/internal/autodb"
	"server_go/internal/dao"
	"server_go/internal/model/entity"
	"server_go/internal/runtime/lock"
	resService "server_go/internal/service/res"
	"server_go/utility/dbcache"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func Login(ctx context.Context, uid int64, loginKey, openid, platform, version string) (g.Map, error) {
	if openid == "" {
		return nil, fmt.Errorf("参数错误: openid 必填")
	}

	lockKey := fmt.Sprintf("user_login:%d", uid)
	token, err := lock.Lock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, fmt.Errorf("系统繁忙，请稍后再试")
	}
	defer func() { _ = lock.Unlock(ctx, lockKey, token) }()

	out := g.Map{"uid": uid}

	var one *entity.User
	err = dao.User.Ctx(ctx).Where("uid", uid).Scan(&one)
	if err != nil {
		return nil, err
	}

	if one != nil {
		if one.Platform != platform || one.Openid != openid {
			return nil, fmt.Errorf("账号信息不匹配")
		}
		out["newbie"] = 0
		out["user"] = one
	} else {
		out["newbie"] = 1
		nowDay := gtime.Now().StartOfDay().Timestamp()
		err = autodb.DB(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, e := dao.User.Ctx(ctx).TX(tx).Data(g.Map{
				"uid": uid, "platform": platform, "openid": openid,
			}).Insert()
			if e != nil {
				return e
			}
			_, e = dao.UserRes.Ctx(ctx).TX(tx).Data(g.Map{
				"uid": uid, "gold": 200, "diamond": 100, "star": 0,
				"tili": 100, "tili_time": 0, "exp": 0, "level": 1, "day_time": nowDay,
			}).Insert()
			return e
		})
		if err != nil {
			return nil, err
		}
		out["user"] = &entity.User{Uid: uint(uid), Platform: platform, Openid: openid}
	}

	bgCtx := autodb.BackgroundWithChannel(ctx)
	go func() {
		defer func() { recover() }()
		_, _ = dao.LogLogin.Ctx(bgCtx).Data(g.Map{"uid": uid, "platform": platform}).Insert()
	}()

	_, err = dao.UserLoginkey.Ctx(ctx).Data(g.Map{
		"uid": uid, "key": loginKey, "ver": version, "time": gtime.Timestamp(),
	}).Save()
	if err != nil {
		return nil, err
	}
	autodb.Redis(ctx).Do(ctx, "SETEX", dbcache.BuildKey(ctx, "login_key", "uid", strconv.FormatInt(uid, 10)), 7200, loginKey)

	out["datas"], err = dao.UserData.Ctx(ctx).Where("uid", uid).All()
	if err != nil {
		return nil, err
	}

	gmVal, err := dao.SysGm.Ctx(ctx).Where("uid", uid).Value("uid")
	if err != nil {
		return nil, err
	}
	if gmVal.IsEmpty() {
		out["gm"] = 0
	} else {
		out["gm"] = 1
	}

	out["items"], err = dao.UserItem.Ctx(ctx).Where("uid", uid).All()
	if err != nil {
		return nil, err
	}

	out["res"], err = resService.Check(ctx, uid)
	if err != nil {
		return nil, err
	}

	out["config"], err = dao.MemConfig.Ctx(ctx).All()
	if err != nil {
		return nil, err
	}

	return out, nil
}
