// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PrfTaskDao 是表 prf_task 的数据访问对象。
type PrfTaskDao struct {
	table    string             // table 是 DAO 底层表名。
	group    string             // group 是当前 DAO 的数据库配置分组名称。
	columns  PrfTaskColumns     // columns 包含表的所有列名，便于使用。
	handlers []gdb.ModelHandler // handlers 用于自定义模型修改。
}

// PrfTaskColumns 定义并存储表 prf_task 的列名。
type PrfTaskColumns struct {
	Id        string //
	Ser       string // 对应工厂的ser
	Tid       string // 任务道具id
	Npc       string //
	StartLoop string // 开启循环的第一个任务有效,每组里如果都是0,不开启循环
}

// prfTaskColumns 保存表 prf_task 的列信息。
var prfTaskColumns = PrfTaskColumns{
	Id:        "id",
	Ser:       "ser",
	Tid:       "tid",
	Npc:       "npc",
	StartLoop: "start_loop",
}

// NewPrfTaskDao 创建并返回用于表数据访问的新 DAO 对象。
func NewPrfTaskDao(handlers ...gdb.ModelHandler) *PrfTaskDao {
	return &PrfTaskDao{
		group:    "default",
		table:    "prf_task",
		columns:  prfTaskColumns,
		handlers: handlers,
	}
}

// DB 获取并返回当前 DAO 底层的原始数据库管理对象。
func (dao *PrfTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 的表名。
func (dao *PrfTaskDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *PrfTaskDao) Columns() PrfTaskColumns {
	return dao.columns
}

// Group 返回当前 DAO 的数据库配置分组名称。
func (dao *PrfTaskDao) Group() string {
	return dao.group
}

// Ctx 为当前 DAO 创建并返回 Model，并自动设置当前操作的上下文。
func (dao *PrfTaskDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction 使用函数 f 包装事务逻辑。
// 如果函数 f 返回非 nil 错误，则回滚事务并返回该错误。
// 如果函数 f 返回 nil，则提交事务并返回 nil。
//
// 注意：不要在函数 f 中提交或回滚事务，
// 因为本函数会自动处理。
func (dao *PrfTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
