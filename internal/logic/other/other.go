package other

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/service"
	"server_go/utility/secretutil"

	"github.com/gogf/gf/v2/frame/g"
)

type sOther struct{}

func init() {
	service.RegisterOther(&sOther{})
}

func (s *sOther) GetResVersion(ctx context.Context, key string) (g.Map, error) {
	return GetResVersion(ctx, key)
}

// GetResVersion checks secret key, prevents replay, returns resource version.
func GetResVersion(ctx context.Context, key string) (g.Map, error) {
	redis := g.Redis()
	rkey := fmt.Sprintf("res_version.%s", key)

	exists, err := redis.Do(ctx, "EXISTS", rkey)
	if err != nil {
		return nil, err
	}
	if exists.Int() > 0 {
		return g.Map{"code": -1036, "msg": "get_res_version: 不能重复调用"}, nil
	}
	_, _ = redis.Do(ctx, "SET", rkey, "1", "EX", 3600)

	if !secretutil.CheckSecret(key) {
		return g.Map{"code": -1, "msg": "参数错误"}, nil
	}

	ver, err := dao.MemConfig.Ctx(ctx).Where("id", 50).Value("value")
	if err != nil {
		return nil, err
	}

	return g.Map{"code": 0, "ver": ver.String()}, nil
}