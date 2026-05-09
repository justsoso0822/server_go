// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserTaskDao is the data access object for the table user_task.
type UserTaskDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserTaskColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserTaskColumns defines and stores column names for the table user_task.
type UserTaskColumns struct {
	Uid    string //
	Taskid string //
	Addtm  string // 添加时间
	Done   string // 是否完成
	Donetm string // 完成时间
}

// userTaskColumns holds the columns for the table user_task.
var userTaskColumns = UserTaskColumns{
	Uid:    "uid",
	Taskid: "taskid",
	Addtm:  "addtm",
	Done:   "done",
	Donetm: "donetm",
}

// NewUserTaskDao creates and returns a new DAO object for table data access.
func NewUserTaskDao(handlers ...gdb.ModelHandler) *UserTaskDao {
	return &UserTaskDao{
		group:    "default",
		table:    "user_task",
		columns:  userTaskColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserTaskDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserTaskDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserTaskDao) Columns() UserTaskColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserTaskDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserTaskDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserTaskDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
