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

// 这段 Lua 在 Redis 内部一次性完成“读取当前值、必要时写入/续期、冲突时拒绝”。
// 扩展知识：Redis 单条 Lua 脚本执行期间是原子的，适合封装这种不能被并发请求打断的状态机。
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
	UID      int64
	LoginKey string
	Platform string
	Version  string
}

type AuthResult struct {
	Code int
	Msg  string
}

func VerifyLoginKey(ctx context.Context, in AuthInput) AuthResult {
	if in.UID == 0 || in.LoginKey == "" || in.Platform == "" || in.Version == "" {
		return AuthResult{Code: -1, Msg: "Verify: 参数错误"}
	}

	rc := autodb.Redis(ctx)
	cacheKey := state.BuildKey(ctx, "login_key", "uid", strconv.FormatInt(in.UID, 10))

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

	// redis.Nil 是 go-redis 表示 key 不存在的哨兵错误，不是真正的系统异常。
	keyData, err := dao.GetUserLoginkey(ctx, in.UID, in.LoginKey)
	if err != nil {
		return AuthResult{Code: -1, Msg: "Verify: 查询失败"}
	}
	if keyData == nil {
		return AuthResult{Code: -1, Msg: "Verify: login_key无效"}
	}

	// Eval 的 KEYS/ARGV 分开传是 Redis Lua 的惯例；集群模式下 key 必须显式放在 KEYS 中，
	// Redis 才能判断脚本访问的槽位。
	result, err := rc.Eval(ctx, loginKeyFillScript, []string{cacheKey}, in.LoginKey, int(loginKeyTTL.Seconds())).Int()
	if err != nil {
		return AuthResult{Code: -1, Msg: "Verify: 缓存写入失败"}
	}
	if result == -1 {
		return AuthResult{Code: -1035, Msg: "Verify: 该账号已在其他地方登陆"}
	}

	return AuthResult{Code: 0}
}
