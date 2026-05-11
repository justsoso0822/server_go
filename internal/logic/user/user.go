package user

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/logic/lock"
	"server_go/internal/model/entity"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

type sUser struct{}

func init() {
	service.RegisterUser(&sUser{})
}

func (s *sUser) Login(ctx context.Context, uid int64, loginKey, openid, platform, version string) (g.Map, error) {
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

	var user *entity.User
	err = dao.User.Ctx(ctx).Where("uid", uid).Scan(&user)
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
		nowDay := gtime.Now().StartOfDay().Timestamp()
		err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
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
		out["user"] = &entity.User{
			Uid: uint(uid), Platform: platform, Openid: openid,
		}
	}

	// 记录登录日志（异步执行并 recover）
	bgCtx := gctx.NeverDone(ctx)
	go func() {
		defer func() { recover() }()
		_, _ = dao.LogLogin.Ctx(bgCtx).Data(g.Map{"uid": uid, "platform": platform}).Insert()
	}()

	// 写入或更新登录密钥。串行化登录后，最后完成的一次登录会成为唯一有效客户端。
	_, err = dao.UserLoginkey.Ctx(ctx).Data(g.Map{
		"uid": uid, "key": loginKey, "ver": version, "time": gtime.Timestamp(),
	}).Save()
	if err != nil {
		return nil, err
	}

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

	out["res"], err = s.GetUserRes(ctx, uid)
	if err != nil {
		return nil, err
	}

	out["config"], err = dao.MemConfig.Ctx(ctx).All()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *sUser) GetUser(ctx context.Context, uid int64) (*entity.User, error) {
	var user *entity.User
	err := dao.User.Ctx(ctx).Where("uid", uid).Scan(&user)
	return user, err
}

func (s *sUser) GetUserRes(ctx context.Context, uid int64) (*entity.UserRes, error) {
	var res *entity.UserRes
	err := dao.UserRes.Ctx(ctx).Where("uid", uid).Scan(&res)
	if err != nil || res == nil {
		return res, err
	}
	nowDay := int(gtime.Now().StartOfDay().Timestamp())
	if res.DayTime != nowDay {
		_, _ = dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{"day_conf": "", "day_time": nowDay}).Update()
		res.DayConf = ""
		res.DayTime = nowDay
	}
	return res, nil
}
