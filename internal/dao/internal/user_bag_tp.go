// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserBagTpDao is the data access object for the table user_bag_tp.
type UserBagTpDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserBagTpColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserBagTpColumns defines and stores column names for the table user_bag_tp.
type UserBagTpColumns struct {
	Id      string //
	Uid     string //
	Chapter string // 副本id - 0 主线
	Time    string // 更新时间
	Itemid  string // 物品id, 0=空
	Count   string // 道具数量
}

// userBagTpColumns holds the columns for the table user_bag_tp.
var userBagTpColumns = UserBagTpColumns{
	Id:      "id",
	Uid:     "uid",
	Chapter: "chapter",
	Time:    "time",
	Itemid:  "itemid",
	Count:   "count",
}

// NewUserBagTpDao creates and returns a new DAO object for table data access.
func NewUserBagTpDao(handlers ...gdb.ModelHandler) *UserBagTpDao {
	return &UserBagTpDao{
		group:    "default",
		table:    "user_bag_tp",
		columns:  userBagTpColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserBagTpDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserBagTpDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserBagTpDao) Columns() UserBagTpColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserBagTpDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserBagTpDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserBagTpDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
