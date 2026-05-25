package res

import (
	"context"

	"server_go/internal/dao"
	"server_go/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func Check(ctx context.Context, uid int64) (*entity.UserRes, error) {
	var res *entity.UserRes
	err := dao.UserRes.Ctx(ctx).Where("uid", uid).Scan(&res)
	if err != nil || res == nil {
		return res, err
	}
	nowDay := int(gtime.Now().StartOfDay().Timestamp())
	if res.DayTime != nowDay {
		_, _ = dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{"day_conf": "", "day_time": nowDay}).Update()
		res.DayConf = ""
		res.DayTime = nowDay
	}
	return res, nil
}
