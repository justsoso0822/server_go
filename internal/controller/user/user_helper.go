package user

import (
	"server_go/api/user/v1"
	"server_go/internal/model"
)

func toUpdateFieldRes(out *model.UpdateFieldOutput) *v1.UpdateFieldRes {
	if out == nil {
		return nil
	}
	return &v1.UpdateFieldRes{
		Res:      out.Res,
		AddValue: out.AddValue,
	}
}
