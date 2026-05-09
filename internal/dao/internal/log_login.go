// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LogLoginDao is the data access object for the table log_login.
type LogLoginDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LogLoginColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LogLoginColumns defines and stores column names for the table log_login.
type LogLoginColumns struct {
	Id       string //
	Uid      string //
	Platform string //
	Time     string //
}

// logLoginColumns holds the columns for the table log_login.
var logLoginColumns = LogLoginColumns{
	Id:       "id",
	Uid:      "uid",
	Platform: "platform",
	Time:     "time",
}

// NewLogLoginDao creates and returns a new DAO object for table data access.
func NewLogLoginDao(handlers ...gdb.ModelHandler) *LogLoginDao {
	return &LogLoginDao{
		group:    "default",
		table:    "log_login",
		columns:  logLoginColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LogLoginDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LogLoginDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LogLoginDao) Columns() LogLoginColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LogLoginDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LogLoginDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LogLoginDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
