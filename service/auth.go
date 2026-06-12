package service

import (
	"context"
	"strconv"
	"time"

	"server_go/dao"
	"server_go/state"
	"server_go/tools/autodb"
)

const (
	loginKeyTTL            = 2 * time.Hour
	loginKeyRenewThreshold = loginKeyTTL / 2
)

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

// 热路径校验 login_key，并且只在剩余 TTL 小于等于阈值时续期。
// 返回值：1=校验通过未续期，2=校验通过且已续期，0=缓存不存在，-1=login_key 冲突。
const loginKeyCheckScript = `
local current = redis.call('GET', KEYS[1])
if current == false then
	return 0
end
if current ~= ARGV[1] then
	return -1
end
local ttl = redis.call('TTL', KEYS[1])
if ttl < 0 or ttl <= tonumber(ARGV[3]) then
	redis.call('EXPIRE', KEYS[1], ARGV[2])
	return 2
end
return 1
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
	ttlSeconds := int(loginKeyTTL.Seconds())
	renewThresholdSeconds := int(loginKeyRenewThreshold.Seconds())

	// 热路径使用 Lua 把 GET、比对、TTL 检查和低频续期合并成一个原子操作。
	// 这样不会每次请求都 EXPIRE，只有剩余有效期低于一半时才续一次。
	checkResult, err := rc.Eval(ctx, loginKeyCheckScript, []string{cacheKey}, in.LoginKey, ttlSeconds, renewThresholdSeconds).Int()
	if err != nil {
		return AuthResult{Code: -1, Msg: "Verify: 缓存读取失败"}
	}
	switch checkResult {
	case 1, 2:
		return AuthResult{Code: 0}
	case -1:
		return AuthResult{Code: -1035, Msg: "Verify: 该账号已在其他地方登陆"}
	}

	// 返回 0 表示 Redis 中没有 login_key，回源 DB 校验后再填充缓存。
	keyData, err := dao.GetUserLoginkey(ctx, in.UID, in.LoginKey)
	if err != nil {
		return AuthResult{Code: -1, Msg: "Verify: 查询失败"}
	}
	if keyData == nil {
		return AuthResult{Code: -1, Msg: "Verify: login_key无效"}
	}

	// Eval 的 KEYS/ARGV 分开传是 Redis Lua 的惯例；集群模式下 key 必须显式放在 KEYS 中，
	// Redis 才能判断脚本访问的槽位。
	fillResult, err := rc.Eval(ctx, loginKeyFillScript, []string{cacheKey}, in.LoginKey, ttlSeconds).Int()
	if err != nil {
		return AuthResult{Code: -1, Msg: "Verify: 缓存写入失败"}
	}
	if fillResult == -1 {
		return AuthResult{Code: -1035, Msg: "Verify: 该账号已在其他地方登陆"}
	}

	return AuthResult{Code: 0}
}
