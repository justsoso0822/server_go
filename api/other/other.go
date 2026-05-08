package other

import (
	"server_go/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type ResVersionReq struct {
	g.Meta `path:"/res_version/:key" method:"get,post" tags:"Other" summary:"获取资源版本号"`
	Key    string `json:"key" in:"path" v:"required"`
}
type ResVersionRes model.ResVersionOutput