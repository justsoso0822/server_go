// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"server_go/internal/autodb"

	"github.com/gogf/gf/v2/database/gdb"
)

// SysDieDao is the data access object for the table sys_die.
type SysDieDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysDieColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysDieColumns defines and stores column names for the table sys_die.
type SysDieColumns struct {
	Uid  string //
	Tips string // 封掉用户登陆时候的错误提示
	Time string //
}

// sysDieColumns holds the columns for the table sys_die.
var sysDieColumns = SysDieColumns{
	Uid:  "uid",
	Tips: "tips",
	Time: "time",
}

// NewSysDieDao creates and returns a new DAO object for table data access.
func NewSysDieDao(handlers ...gdb.ModelHandler) *SysDieDao {
	return &SysDieDao{
		group:    "default",
		table:    "sys_die",
		columns:  sysDieColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the raw database management object of the current DAO using request context.
func (dao *SysDieDao) DB(ctx context.Context) gdb.DB {
	return autodb.DB(ctx, dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysDieDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysDieDao) Columns() SysDieColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysDieDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysDieDao) Ctx(ctx context.Context) *gdb.Model {
	model := autodb.DB(ctx, dao.group).Model(dao.table)
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
func (dao *SysDieDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
