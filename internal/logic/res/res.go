package res

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/logic/gamelog"
	"server_go/internal/logic/lock"
	"server_go/internal/model/entity"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sRes struct{}

func init() {
	service.RegisterRes(&sRes{})
}

func (s *sRes) UpdateDiamond(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "diamond", "钻石")
}

func (s *sRes) UpdateGold(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "gold", "金币")
}

func (s *sRes) UpdateTili(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "tili", "体力")
}

func (s *sRes) UpdateExp(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "exp", "经验")
}

func (s *sRes) UpdateStar(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "star", "星星")
}

func updateResField(ctx context.Context, uid int64, cnt int64, reason string, field, resName string) (g.Map, error) {
	lockKey := fmt.Sprintf("update_%s:%d", field, uid)
	token, err := lock.Lock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, fmt.Errorf("系统繁忙，请稍后再试")
	}
	defer func() { _ = lock.Unlock(ctx, lockKey, token) }()

	var res *entity.UserRes
	err = dao.UserRes.Ctx(ctx).Where("uid", uid).Scan(&res)
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

	newCnt := oldCnt + cnt
	if newCnt < 0 {
		newCnt = 0
	}
	if newCnt == oldCnt {
		return g.Map{"res": res, "add_value": 0}, nil
	}

	_, err = dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{field: newCnt}).Update()
	if err != nil {
		gamelog.Log(ctx, uid, fmt.Sprintf("更新用户资源失败 %s %d %s %v", field, cnt, reason, err))
		return nil, err
	}

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

	gamelog.TraceRes(ctx, uid, oldCnt, newCnt, resName, reason)
	return g.Map{"res": res, "add_value": newCnt - oldCnt}, nil
}
