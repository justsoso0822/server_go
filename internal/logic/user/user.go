package user

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"server_go/internal/consts"
	"server_go/internal/dao"
	"server_go/internal/logic/gamelog"
	"server_go/internal/logic/lock"
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

func (s *sUser) UpdateDiamond(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "diamond", "钻石")
}

func (s *sUser) UpdateGold(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "gold", "金币")
}

func (s *sUser) UpdateTili(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "tili", "体力")
}

func (s *sUser) UpdateExp(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "exp", "经验")
}

func (s *sUser) UpdateStar(ctx context.Context, in *model.UpdateFieldInput) (*model.UpdateFieldOutput, error) {
	return updateResField(ctx, in, "star", "星星")
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

func updateResField(ctx context.Context, in *model.UpdateFieldInput, field, resName string) (*model.UpdateFieldOutput, error) {
	lockKey := fmt.Sprintf("update_%s:%d", field, in.Uid)
	token, err := lock.Lock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, fmt.Errorf("系统繁忙，请稍后再试")
	}
	defer func() { _ = lock.Unlock(ctx, lockKey, token) }()

	var res *entity.UserRes
	err = dao.UserRes.Ctx(ctx).Where("uid", in.Uid).Scan(&res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("用户资源不存在")
	}

	var oldCnt int64
	switch field {
	case "diamond":
		oldCnt = int64(res.Diamond)
	case "gold":
		oldCnt = int64(res.Gold)
	case "tili":
		oldCnt = int64(res.Tili)
	case "exp":
		oldCnt = int64(res.Exp)
	case "star":
		oldCnt = int64(res.Star)
	}

	newCnt := oldCnt + in.Cnt
	if newCnt < 0 {
		newCnt = 0
	}
	if newCnt == oldCnt {
		return &model.UpdateFieldOutput{Res: res, AddValue: 0}, nil
	}

	_, err = dao.UserRes.Ctx(ctx).Where("uid", in.Uid).Data(g.Map{field: newCnt}).Update()
	if err != nil {
		gamelog.Log(ctx, in.Uid, fmt.Sprintf("更新用户资源失败 %s %d %s %v", field, in.Cnt, in.Reason, err))
		return nil, err
	}

	// 更新结构体
	switch field {
	case "diamond":
		res.Diamond = int(newCnt)
	case "gold":
		res.Gold = int(newCnt)
	case "tili":
		res.Tili = int(newCnt)
	case "exp":
		res.Exp = int(newCnt)
	case "star":
		res.Star = int(newCnt)
	}

	gamelog.TraceRes(ctx, in.Uid, oldCnt, newCnt, resName, in.Reason)
	return &model.UpdateFieldOutput{Res: res, AddValue: newCnt - oldCnt}, nil
}

// --- 工具方法 ---

func ParseRes(items interface{}) []consts.ResItem {
	switch v := items.(type) {
	case []consts.ResItem:
		return v
	case string:
		return parseResString(v)
	default:
		return nil
	}
}

func parseResString(s string) []consts.ResItem {
	nums := PickNumbers(s)
	if len(nums) == 0 || len(nums)%3 != 0 {
		return nil
	}
	result := make([]consts.ResItem, 0, len(nums)/3)
	for i := 0; i < len(nums); i += 3 {
		result = append(result, consts.ResItem{Type: nums[i], Id: nums[i+1], Cnt: nums[i+2]})
	}
	return result
}

// PickNumbers 从字符串中提取所有整数。
func PickNumbers(s string) []int {
	var result []int
	current := ""
	for _, c := range s {
		if (c >= '0' && c <= '9') || c == '-' || c == '+' || c == '.' {
			current += string(c)
		} else {
			if current != "" {
				if n, err := strconv.Atoi(current); err == nil {
					result = append(result, n)
				} else if f, err := strconv.ParseFloat(current, 64); err == nil {
					result = append(result, int(math.Floor(f)))
				}
				current = ""
			}
		}
	}
	if current != "" {
		if n, err := strconv.Atoi(current); err == nil {
			result = append(result, n)
		} else if f, err := strconv.ParseFloat(current, 64); err == nil {
			result = append(result, int(math.Floor(f)))
		}
	}
	return result
}

func init() {
	_ = strings.Join
}
