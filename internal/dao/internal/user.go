// ==========================================================================
// 代码由 GoFrame CLI 工具生成并维护。请勿编辑。
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserDao 是表 user 的数据访问对象。
type UserDao struct {
	table    string             // table 是 DAO 底层表名。
	group    string             // group 是当前 DAO 的数据库配置分组名称。
	columns  UserColumns        // columns 包含表的所有列名，便于使用。
	handlers []gdb.ModelHandler // handlers 用于自定义模型修改。
}

// UserColumns 定义并存储表 user 的列名。
type UserColumns struct {
	Uid      string //
	Platform string //
	Openid   string //
	Sid      string //
	Name     string //
	Head     string //
	Born     string //
}

// userColumns 保存表 user 的列信息。
var userColumns = UserColumns{
	Uid:      "uid",
	Platform: "platform",
	Openid:   "openid",
	Sid:      "sid",
	Name:     "name",
	Head:     "head",
	Born:     "born",
}

// NewUserDao 创建并返回用于表数据访问的新 DAO 对象。
func NewUserDao(handlers ...gdb.ModelHandler) *UserDao {
	return &UserDao{
		group:    "default",
		table:    "user",
		columns:  userColumns,
		handlers: handlers,
	}
}

// DB 获取并返回当前 DAO 底层的原始数据库管理对象。
func (dao *UserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 的表名。
func (dao *UserDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *UserDao) Columns() UserColumns {
	return dao.columns
}

// Group 返回当前 DAO 的数据库配置分组名称。
func (dao *UserDao) Group() string {
	return dao.group
}

// Ctx 为当前 DAO 创建并返回 Model，并自动设置当前操作的上下文。
func (dao *UserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
