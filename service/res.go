package service

import (
	"context"
	"fmt"

	"server_gin/dao"
	"server_gin/dao/model"
	"server_gin/runtime"
)

func UpdateDiamond(ctx context.Context, uid int64, cnt int64, reason string) (map[string]any, error) {
	return updateResField(ctx, uid, cnt, reason, dao.UserResFieldDiamond, "钻石")
}

func UpdateGold(ctx context.Context, uid int64, cnt int64, reason string) (map[string]any, error) {
	return updateResField(ctx, uid, cnt, reason, dao.UserResFieldGold, "金币")
}

func UpdateTili(ctx context.Context, uid int64, cnt int64, reason string) (map[string]any, error) {
	return updateResField(ctx, uid, cnt, reason, dao.UserResFieldTili, "体力")
}

func UpdateExp(ctx context.Context, uid int64, cnt int64, reason string) (map[string]any, error) {
	return updateResField(ctx, uid, cnt, reason, dao.UserResFieldExp, "经验")
}

func UpdateStar(ctx context.Context, uid int64, cnt int64, reason string) (map[string]any, error) {
	return updateResField(ctx, uid, cnt, reason, dao.UserResFieldStar, "星星")
}

func updateResField(ctx context.Context, uid int64, cnt int64, reason string, field dao.UserResField, resName string) (map[string]any, error) {
	lockKey := fmt.Sprintf("update_%s:%d", field, uid)
	token, err := runtime.Lock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, fmt.Errorf("系统繁忙，请稍后再试")
	}
	defer func() { _ = runtime.Unlock(ctx, lockKey, token) }()

	res, err := dao.GetUserRes(ctx, uid)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("用户资源不存在")
	}

	var oldCnt int64
	switch field {
	case dao.UserResFieldDiamond:
		oldCnt = int64(res.Diamond)
	case dao.UserResFieldGold:
		oldCnt = int64(res.Gold)
	case dao.UserResFieldTili:
		oldCnt = int64(res.Tili)
	case dao.UserResFieldExp:
		oldCnt = int64(res.Exp)
	case dao.UserResFieldStar:
		oldCnt = int64(res.Star)
	}

	newCnt := oldCnt + cnt
	if cnt < 0 && newCnt < 0 {
		return nil, fmt.Errorf("%s余额不足", resName)
	}
	if newCnt < 0 {
		newCnt = 0
	}
	if newCnt == oldCnt {
		return map[string]any{"res": res, "add_value": int64(0)}, nil
	}

	if err = dao.UpdateUserResField(ctx, uid, field, newCnt); err != nil {
		runtime.LogMsg(ctx, uid, fmt.Sprintf("更新用户资源失败 %s %d %s %v", field, cnt, reason, err))
		return nil, err
	}

	updateResStruct(res, field, int(newCnt))
	runtime.TraceRes(ctx, uid, oldCnt, newCnt, resName, reason)

	return map[string]any{"res": res, "add_value": newCnt - oldCnt}, nil
}

func updateResStruct(res *model.UserRes, field dao.UserResField, val int) {
	v := int32(val)
	switch field {
	case dao.UserResFieldDiamond:
		res.Diamond = v
	case dao.UserResFieldGold:
		res.Gold = v
	case dao.UserResFieldTili:
		res.Tili = v
	case dao.UserResFieldExp:
		res.Exp = v
	case dao.UserResFieldStar:
		res.Star = v
	}
}
