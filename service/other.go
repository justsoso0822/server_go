package service

import (
	"context"
	"time"

	"server_go/dao"
	"server_go/state"
	"server_go/tools/autodb"
	secretutil "server_go/tools/secret"
)

func GetResVersion(ctx context.Context, key string) (map[string]interface{}, error) {
	if !secretutil.CheckSecret(key) {
		return map[string]interface{}{"code": -1, "msg": "参数错误"}, nil
	}

	rc := autodb.Redis(ctx)
	rkey := state.BuildKey(ctx, "res_version", key)

	ok, err := rc.SetNX(ctx, rkey, "1", time.Hour).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return map[string]interface{}{"code": -1036, "msg": "get_res_version: 不能重复调用"}, nil
	}

	ver, err := dao.GetMemConfigValue(ctx, 50)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"code": 0, "ver": ver}, nil
}
