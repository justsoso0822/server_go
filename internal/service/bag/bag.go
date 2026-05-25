package bag

import (
	"context"

	"server_go/internal/dao"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func GetUserBag(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	rows, err := dao.UserBag.Ctx(ctx).Where("uid", uid).Where("chapter", chapter).All()
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = gdb.Result{}
	}
	return g.Map{"uid": uid, "chapter": chapter, "bag": rows}, nil
}

func GetUserBagTp(ctx context.Context, uid int64, chapter int) (g.Map, error) {
	rows, err := dao.UserBagTp.Ctx(ctx).Where("uid", uid).Where("chapter", chapter).All()
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = gdb.Result{}
	}
	return g.Map{"uid": uid, "chapter": chapter, "bag": rows}, nil
}
