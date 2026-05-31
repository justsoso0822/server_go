package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"server_go/dao"
	"server_go/state"
	"server_go/tools/autodb"

	"github.com/redis/go-redis/v9"
)

const loginKeyTTL = 2 * time.Hour

const loginKeyFillScript = `
local current = redis.call('GET', KEYS[1])
if current == false then
	redis.call('SETEX', KEYS[1], ARGV[2], ARGV[1])
	return 1
end
if current == ARGV[1] then
	redis.call('EXPIRE', KEYS[1], ARGV[2])
	return 1
end
return -1
`

type AuthInput struct {
	Uid      int64
	LoginKey string
	Platform string
	Version  string
}

type AuthResult struct {
	Code int
	Msg  string
}

func VerifyLoginKey(ctx context.Context, in AuthInput) AuthResult {
	if in.Uid == 0 || in.LoginKey == "" || in.Platform == "" || in.Version == "" {
		return AuthResult{Code: -1, Msg: "Verify: 参数错误"}
	}

	rc := autodb.Redis(ctx)
	cacheKey := state.BuildKey(ctx, "login_key", "uid", strconv.FormatInt(in.Uid, 10))

	cached, err := rc.Get(ctx, cacheKey).Result()
	switch {
	case err == nil:
		if cached == in.LoginKey {
			return AuthResult{Code: 0}
		}
		return AuthResult{Code: -1035, Msg: "Verify: 该账号已在其他地方登陆"}
	case !errors.Is(err, redis.Nil):
		return AuthResult{Code: -1, Msg: "Verify: 缓存读取失败"}
	}

	keyData, err := dao.GetUserLoginkey(ctx, in.Uid, in.LoginKey)
	if err != nil {
		return AuthResult{Code: -1, Msg: "Verify: 查询失败"}
	}
	if keyData == nil {
		return AuthResult{Code: -1, Msg: "Verify: login_key无效"}
	}

	result, err := rc.Eval(ctx, loginKeyFillScript, []string{cacheKey}, in.LoginKey, int(loginKeyTTL.Seconds())).Int()
	if err != nil {
		return AuthResult{Code: -1, Msg: "Verify: 缓存写入失败"}
	}
	if result == -1 {
		return AuthResult{Code: -1035, Msg: "Verify: 该账号已在其他地方登陆"}
	}

	return AuthResult{Code: 0}
}
