package game

import (
	"context"

	"server_go/internal/dao"
	"server_go/internal/model"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sGame struct{}

func init() {
	service.RegisterGame(&sGame{})
}

func (s *sGame) Online(ctx context.Context, in *model.OnlineInput) error {
	now := gtime.Now()
	dayStr := now.Format("Y-m-d, H:00:00")

	row, err := dao.UserOnline.Ctx(ctx).Where("uid", in.Uid).Where("day", dayStr).One()
	if err != nil {
		return err
	}

	nowTime := now.Format("Y-m-d H:i:s")
	seconds := in.Seconds
	if !row.IsEmpty() {
		seconds += row["tm_online"].Int64()
		_, err = dao.UserOnline.Ctx(ctx).Where("uid", in.Uid).Where("day", dayStr).
			Data(g.Map{"tm_online": seconds, "tm_update": nowTime}).Update()
	} else {
		_, err = dao.UserOnline.Ctx(ctx).Data(g.Map{
			"uid": in.Uid, "day": dayStr, "tm_online": seconds, "tm_update": nowTime,
		}).Insert()
	}
	return err
}
