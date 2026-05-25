package res

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/model/entity"
	"server_go/internal/runtime/gamelog"
	"server_go/internal/runtime/lock"

	"github.com/gogf/gf/v2/frame/g"
)

func UpdateDiamond(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "diamond", "钻石")
}

func UpdateGold(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "gold", "金币")
}

func UpdateTili(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "tili", "体力")
}

func UpdateExp(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
	return updateResField(ctx, uid, cnt, reason, "exp", "经验")
}

func UpdateStar(ctx context.Context, uid int64, cnt int64, reason string) (g.Map, error) {
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

	var one *entity.UserRes
	err = dao.UserRes.Ctx(ctx).Where("uid", uid).Scan(&one)
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, fmt.Errorf("用户资源不存在")
	}

	var oldCnt int64
	switch field {
	case "diamond":
		oldCnt = int64(one.Diamond)
	case "gold":
		oldCnt = int64(one.Gold)
	case "tili":
		oldCnt = int64(one.Tili)
	case "exp":
		oldCnt = int64(one.Exp)
	case "star":
		oldCnt = int64(one.Star)
	}

	newCnt := oldCnt + cnt
	if cnt < 0 && newCnt < 0 {
		return nil, fmt.Errorf("%s余额不足", resName)
	}
	if newCnt < 0 {
		newCnt = 0
	}
	if newCnt == oldCnt {
		return g.Map{"res": one, "add_value": 0}, nil
	}

	_, err = dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{field: newCnt}).Update()
	if err != nil {
		gamelog.Log(ctx, uid, fmt.Sprintf("更新用户资源失败 %s %d %s %v", field, cnt, reason, err))
		return nil, err
	}

	switch field {
	case "diamond":
		one.Diamond = int(newCnt)
	case "gold":
		one.Gold = int(newCnt)
	case "tili":
		one.Tili = int(newCnt)
	case "exp":
		one.Exp = int(newCnt)
	case "star":
		one.Star = int(newCnt)
	}

	gamelog.TraceRes(ctx, uid, oldCnt, newCnt, resName, reason)
	return g.Map{"res": one, "add_value": newCnt - oldCnt}, nil
}
