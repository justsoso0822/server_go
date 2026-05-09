// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserBagDao is the data access object for the table user_bag.
type UserBagDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserBagColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserBagColumns defines and stores column names for the table user_bag.
type UserBagColumns struct {
	Id      string //
	Uid     string //
	Chapter string // 副本id - 0-主格子
	Time    string // 操作更新时间
	Itemid  string // 物品id, 0=空
	Info    string // 格子信息
	Type    string // 格子类型,0-普通,1-道具购买
}

// userBagColumns holds the columns for the table user_bag.
var userBagColumns = UserBagColumns{
	Id:      "id",
	Uid:     "uid",
	Chapter: "chapter",
	Time:    "time",
	Itemid:  "itemid",
	Info:    "info",
	Type:    "type",
}

// NewUserBagDao creates and returns a new DAO object for table data access.
func NewUserBagDao(handlers ...gdb.ModelHandler) *UserBagDao {
	return &UserBagDao{
		group:    "default",
		table:    "user_bag",
		columns:  userBagColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserBagDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserBagDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserBagDao) Columns() UserBagColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserBagDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserBagDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserBagDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
