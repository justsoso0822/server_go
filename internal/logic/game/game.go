package game

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sGame struct{}

func init() {
	service.RegisterGame(&sGame{})
}

func (s *sGame) Online(ctx context.Context, uid, seconds int64) (map[string]interface{}, error) {
	return Online(ctx, uid, seconds)
}
func (s *sGame) ServerTime() map[string]interface{} {
	return ServerTime()
}

// Online records user online duration by hour.
func Online(ctx context.Context, uid, seconds int64) (g.Map, error) {
	now := gtime.Now()
	dayStr := fmt.Sprintf("%d-%02d-%02d, %02d:00:00", now.Year(), now.Month(), now.Day(), now.Hour())

	row, err := dao.UserOnline.Ctx(ctx).Where("uid", uid).Where("day", dayStr).One()
	if err != nil {
		return nil, err
	}

	nowTime := gtime.Now().Format("Y-m-d H:i:s")
	if !row.IsEmpty() {
		seconds += row["tm_online"].Int64()
		_, err = dao.UserOnline.Ctx(ctx).Where("uid", uid).Where("day", dayStr).
			Data(g.Map{"tm_online": seconds, "tm_update": nowTime}).Update()
	} else {
		_, err = dao.UserOnline.Ctx(ctx).Data(g.Map{
			"uid": uid, "day": dayStr, "tm_online": seconds, "tm_update": nowTime,
		}).Insert()
	}
	if err != nil {
		return nil, err
	}

	return g.Map{"code": 0, "now": gtime.TimestampMilli()}, nil
}

// ServerTime returns the current server time.
func ServerTime() g.Map {
	return g.Map{"code": 0, "now": gtime.TimestampMilli()}
}