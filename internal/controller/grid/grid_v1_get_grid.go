package grid

import (
	"context"

	"server_go/api/bag/v1"
	gridV1 "server_go/api/grid/v1"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) GetGrid(ctx context.Context, req *gridV1.GetGridReq) (res *gridV1.GetGridRes, err error) {
	out, err := service.Grid().GetGrid(ctx, req.Uid, req.Chapter)
	if err != nil {
		return nil, err
	}
	
	var bag, bagTp *v1.BagRes
	if out["bag"] != nil {
		bagData := out["bag"].(g.Map)
		bag = &v1.BagRes{
			Uid:     bagData["uid"].(int64),
			Chapter: bagData["chapter"].(int),
			Bag:     bagData["bag"].(gdb.Result),
		}
	}
	if out["bag_tp"] != nil {
		bagTpData := out["bag_tp"].(g.Map)
		bagTp = &v1.BagRes{
			Uid:     bagTpData["uid"].(int64),
			Chapter: bagTpData["chapter"].(int),
			Bag:     bagTpData["bag"].(gdb.Result),
		}
	}
	
	return &gridV1.GetGridRes{
		Bag:   bag,
		BagTp: bagTp,
		Tasks: out["tasks"],
	}, nil
}
