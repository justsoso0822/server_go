package service

import (
	"context"
	"time"

	"server_gin/autodb"
	"server_gin/dao"
	"server_gin/runtime"
	"server_gin/secretutil"
)

func GetResVersion(ctx context.Context, key string) (map[string]interface{}, error) {
	rc := autodb.Redis(ctx)
	rkey := runtime.BuildKey(ctx, "res_version", key)

	exists, err := rc.Exists(ctx, rkey).Result()
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return map[string]interface{}{"code": -1036, "msg": "get_res_version: 不能重复调用"}, nil
	}
	rc.Set(ctx, rkey, "1", time.Hour)

	if !secretutil.CheckSecret(key) {
		return map[string]interface{}{"code": -1, "msg": "参数错误"}, nil
	}

	ver, err := dao.GetMemConfigValue(ctx, 50)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"code": 0, "ver": ver}, nil
}
