// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PrfTaskDao is the data access object for the table prf_task.
type PrfTaskDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PrfTaskColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PrfTaskColumns defines and stores column names for the table prf_task.
type PrfTaskColumns struct {
	Id        string //
	Ser       string // 对应工厂的ser
	Tid       string // 任务道具id
	Npc       string //
	StartLoop string // 开启循环的第一个任务有效,每组里如果都是0,不开启循环
}

// prfTaskColumns holds the columns for the table prf_task.
var prfTaskColumns = PrfTaskColumns{
	Id:        "id",
	Ser:       "ser",
	Tid:       "tid",
	Npc:       "npc",
	StartLoop: "start_loop",
}

// NewPrfTaskDao creates and returns a new DAO object for table data access.
func NewPrfTaskDao(handlers ...gdb.ModelHandler) *PrfTaskDao {
	return &PrfTaskDao{
		group:    "default",
		table:    "prf_task",
		columns:  prfTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PrfTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PrfTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PrfTaskDao) Columns() PrfTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PrfTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PrfTaskDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *PrfTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
