// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"server_go/internal/autodb"

	"github.com/gogf/gf/v2/database/gdb"
)

// PrfResDao is the data access object for the table prf_res.
type PrfResDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PrfResColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PrfResColumns defines and stores column names for the table prf_res.
type PrfResColumns struct {
	Id   string // 资源ID
	Name string // 资源名称
	Tips string // 资源说明,用法提示
}

// prfResColumns holds the columns for the table prf_res.
var prfResColumns = PrfResColumns{
	Id:   "id",
	Name: "name",
	Tips: "tips",
}

// NewPrfResDao creates and returns a new DAO object for table data access.
func NewPrfResDao(handlers ...gdb.ModelHandler) *PrfResDao {
	return &PrfResDao{
		group:    "default",
		table:    "prf_res",
		columns:  prfResColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the raw database management object of the current DAO using request context.
func (dao *PrfResDao) DB(ctx context.Context) gdb.DB {
	return autodb.DB(ctx, dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PrfResDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PrfResDao) Columns() PrfResColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PrfResDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PrfResDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PrfResDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
