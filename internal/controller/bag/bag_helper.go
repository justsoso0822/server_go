package bag

import (
	"server_go/api/bag/v1"
	"server_go/internal/model"
)

func ToBagRes(out *model.BagOutput) *v1.BagRes {
	if out == nil {
		return nil
	}
	return &v1.BagRes{
		Uid:     out.Uid,
		Chapter: out.Chapter,
		Bag:     out.Bag,
	}
}
