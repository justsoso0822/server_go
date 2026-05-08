package bag

import (
	"context"

	"server_go/internal/dao"
	"server_go/internal/model"
	"server_go/internal/service"
)

type sBag struct{}

func init() {
	service.RegisterBag(&sBag{})
}

func (s *sBag) GetUserBag(ctx context.Context, in *model.BagInput) (*model.BagOutput, error) {
	rows, err := dao.UserBag.Ctx(ctx).Where("uid", in.Uid).Where("chapter", in.Chapter).All()
	if err != nil {
		return nil, err
	}
	return &model.BagOutput{Uid: in.Uid, Chapter: in.Chapter, Bag: rows}, nil
}

func (s *sBag) GetUserBagTp(ctx context.Context, in *model.BagInput) (*model.BagOutput, error) {
	rows, err := dao.UserBagTp.Ctx(ctx).Where("uid", in.Uid).Where("chapter", in.Chapter).All()
	if err != nil {
		return nil, err
	}
	return &model.BagOutput{Uid: in.Uid, Chapter: in.Chapter, Bag: rows}, nil
}