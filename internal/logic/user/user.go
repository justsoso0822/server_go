package user

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/model"
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

func (s *sUser) Login(ctx context.Context, in *model.LoginInput) (*model.LoginOutput, error) {
	if in.Openid == "" {
		return nil, fmt.Errorf("参数错误: openid 必填")
	}

	out := &model.LoginOutput{Uid: in.Uid}

	var user *entity.User
	err := dao.User.Ctx(ctx).Where("uid", in.Uid).Scan(&user)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if user.Platform != in.Platform || user.Openid != in.Openid {
			return nil, fmt.Errorf("账号信息不匹配")
		}
		out.Newbie = 0
		out.User = user
	} else {
		out.Newbie = 1
		nowDay := gtime.Now().StartOfDay().Timestamp()
		err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, e := dao.User.Ctx(ctx).TX(tx).Data(g.Map{
				"uid": in.Uid, "platform": in.Platform, "openid": in.Openid,
			}).Insert()
			if e != nil {
				return e
			}
			_, e = dao.UserRes.Ctx(ctx).TX(tx).Data(g.Map{
				"uid": in.Uid, "gold": 200, "diamond": 100, "star": 0,
				"tili": 100, "tili_time": 0, "exp": 0, "level": 1, "day_time": nowDay,
			}).Insert()
			return e
		})
		if err != nil {
			return nil, err
		}
		out.User = &entity.User{
			Uid: uint(in.Uid), Platform: in.Platform, Openid: in.Openid,
		}
	}

	// 记录登录日志（异步执行并 recover）
	bgCtx := gctx.NeverDone(ctx)
	go func() {
		defer func() { recover() }()
		_, _ = dao.LogLogin.Ctx(bgCtx).Data(g.Map{"uid": in.Uid, "platform": in.Platform}).Insert()
	}()

	// 写入或更新登录密钥
	_, err = dao.UserLoginkey.Ctx(ctx).Data(g.Map{
		"uid": in.Uid, "key": in.LoginKey, "ver": in.Version, "time": gtime.Timestamp(),
	}).Save()
	if err != nil {
		return nil, err
	}

	out.Datas, err = dao.UserData.Ctx(ctx).Where("uid", in.Uid).All()
	if err != nil {
		return nil, err
	}

	gmVal, err := dao.SysGm.Ctx(ctx).Where("uid", in.Uid).Value("uid")
	if err != nil {
		return nil, err
	}
	if gmVal.IsEmpty() {
		out.Gm = 0
	} else {
		out.Gm = 1
	}

	out.Items, err = dao.UserItem.Ctx(ctx).Where("uid", in.Uid).All()
	if err != nil {
		return nil, err
	}

	out.Res, err = s.GetUserRes(ctx, in.Uid)
	if err != nil {
		return nil, err
	}

	out.Config, err = dao.MemConfig.Ctx(ctx).All()
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
