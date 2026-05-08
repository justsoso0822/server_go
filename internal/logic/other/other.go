package other

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/utility/secretutil"

	"github.com/gogf/gf/v2/frame/g"
)

type sOther struct{}

func init() {
	service.RegisterOther(&sOther{})
}

func (s *sOther) GetResVersion(ctx context.Context, key string) (*model.ResVersionOutput, error) {
	redis := g.Redis()
	rkey := fmt.Sprintf("res_version.%s", key)

	exists, err := redis.Do(ctx, "EXISTS", rkey)
	if err != nil {
		return nil, err
	}
	if exists.Int() > 0 {
		return &model.ResVersionOutput{Code: -1036, Msg: "get_res_version: 不能重复调用"}, nil
	}
	_, _ = redis.Do(ctx, "SET", rkey, "1", "EX", 3600)

	if !secretutil.CheckSecret(key) {
		return &model.ResVersionOutput{Code: -1, Msg: "参数错误"}, nil
	}

	ver, err := dao.MemConfig.Ctx(ctx).Where("id", 50).Value("value")
	if err != nil {
		return nil, err
	}

	return &model.ResVersionOutput{Code: 0, Ver: ver.String()}, nil
}