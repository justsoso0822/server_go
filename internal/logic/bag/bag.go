package bag

import (
	"context"

	"server_go/internal/dao"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sBag struct{}

func init() {
	service.RegisterBag(&sBag{})
}

func (s *sBag) GetUserBag(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	rows, err := dao.UserBag.Ctx(ctx).Where("uid", uid).Where("chapter", chapter).All()
	if err != nil {
		return nil, err
	}
	return g.Map{"uid": uid, "chapter": chapter, "bag": rows}, nil
}

func (s *sBag) GetUserBagTp(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	rows, err := dao.UserBagTp.Ctx(ctx).Where("uid", uid).Where("chapter", chapter).All()
	if err != nil {
		return nil, err
	}
	return g.Map{"uid": uid, "chapter": chapter, "bag": rows}, nil
}
