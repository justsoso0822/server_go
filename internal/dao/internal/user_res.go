// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserResDao is the data access object for the table user_res.
type UserResDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserResColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserResColumns defines and stores column names for the table user_res.
type UserResColumns struct {
	Uid      string //
	Gold     string //
	Diamond  string //
	Star     string //
	Tili     string //
	TiliTime string //
	Exp      string //
	Level    string //
	DayConf  string // 每日重置的数据
	DayTime  string // 上次重置时间
}

// userResColumns holds the columns for the table user_res.
var userResColumns = UserResColumns{
	Uid:      "uid",
	Gold:     "gold",
	Diamond:  "diamond",
	Star:     "star",
	Tili:     "tili",
	TiliTime: "tili_time",
	Exp:      "exp",
	Level:    "level",
	DayConf:  "day_conf",
	DayTime:  "day_time",
}

// NewUserResDao creates and returns a new DAO object for table data access.
func NewUserResDao(handlers ...gdb.ModelHandler) *UserResDao {
	return &UserResDao{
		group:    "default",
		table:    "user_res",
		columns:  userResColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserResDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserResDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserResDao) Columns() UserResColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserResDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserResDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserResDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
