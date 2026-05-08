// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserOnlineDao is the data access object for the table user_online.
type UserOnlineDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserOnlineColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserOnlineColumns defines and stores column names for the table user_online.
type UserOnlineColumns struct {
	Uid      string //
	Day      string //
	TmOnline string //
	TmUpdate string //
}

// userOnlineColumns holds the columns for the table user_online.
var userOnlineColumns = UserOnlineColumns{
	Uid:      "uid",
	Day:      "day",
	TmOnline: "tm_online",
	TmUpdate: "tm_update",
}

// NewUserOnlineDao creates and returns a new DAO object for table data access.
func NewUserOnlineDao(handlers ...gdb.ModelHandler) *UserOnlineDao {
	return &UserOnlineDao{
		group:    "default",
		table:    "user_online",
		columns:  userOnlineColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserOnlineDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserOnlineDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserOnlineDao) Columns() UserOnlineColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserOnlineDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserOnlineDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserOnlineDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
