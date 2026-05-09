// =================================================================================
// 此文件由 GoFrame CLI 工具自动生成，可按需修改。
// =================================================================================

package dao

import (
	"server_go/internal/dao/internal"
)

// prfTaskDao 是表 prf_task 的数据访问对象。
// 可按需在此定义自定义方法以扩展功能。
type prfTaskDao struct {
	*internal.PrfTaskDao
}

var (
	// PrfTask 是表 prf_task 操作的全局访问对象。
	PrfTask = prfTaskDao{internal.NewPrfTaskDao()}
)

// 在下方添加自定义方法和功能。
