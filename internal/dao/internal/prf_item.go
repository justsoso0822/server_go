// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"server_go/internal/autodb"

	"github.com/gogf/gf/v2/database/gdb"
)

// PrfItemDao is the data access object for the table prf_item.
type PrfItemDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PrfItemColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PrfItemColumns defines and stores column names for the table prf_item.
type PrfItemColumns struct {
	Id         string //
	Ser        string // 道具系列
	Type       string // 道具类型  0 普通道具  1 工厂道具  2 转换器  3 点击使用的道具  4 拖拽到目标使用的道具
	Cost       string // 工厂类型道具产出消耗的体力
	Lv         string // 道具等级
	From       string // 来源id
	Name       string // 道具名称
	Tips       string // 道具介绍
	Star       string // 该道具对应的绿星数量
	Next       string // 合成道具id
	AutoDrop   string // 工厂是否自动产出道具
	NeedOpen   string // 是否需要等待开启
	CoolDrop   string // 一次CD掉落个数
	CoolTime   string // 冷却时间
	CoolMaxnum string // 冷却时间累积次数
	CoolMoney  string // 冷却加速花费钻石（最大值）
	Sell       string // 出售价格-金币
	Use        string // 可直接使用的道具，boxid
	Die        string // 特殊-掉落完后变成什么  <0 直接消失  ==0 进入cd  >0 变成指定道具, 如果道具不存在也消失
	Exp        string // 该等级道具合成下一级道具时是否会产出经验
	Gold       string // 副本合成掉落金币
	Rare       string // 0=普通，1=稀有（出售时跳出提醒框）
	Tids       string // 转换器内可投入物品的id
	Count      string // 转换器内对应投入物品的数量
}

// prfItemColumns holds the columns for the table prf_item.
var prfItemColumns = PrfItemColumns{
	Id:         "id",
	Ser:        "ser",
	Type:       "type",
	Cost:       "cost",
	Lv:         "lv",
	From:       "from",
	Name:       "name",
	Tips:       "tips",
	Star:       "star",
	Next:       "next",
	AutoDrop:   "auto_drop",
	NeedOpen:   "need_open",
	CoolDrop:   "cool_drop",
	CoolTime:   "cool_time",
	CoolMaxnum: "cool_maxnum",
	CoolMoney:  "cool_money",
	Sell:       "sell",
	Use:        "use",
	Die:        "die",
	Exp:        "exp",
	Gold:       "gold",
	Rare:       "rare",
	Tids:       "tids",
	Count:      "count",
}

// NewPrfItemDao creates and returns a new DAO object for table data access.
func NewPrfItemDao(handlers ...gdb.ModelHandler) *PrfItemDao {
	return &PrfItemDao{
		group:    "default",
		table:    "prf_item",
		columns:  prfItemColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the raw database management object of the current DAO using request context.
func (dao *PrfItemDao) DB(ctx context.Context) gdb.DB {
	return autodb.DB(ctx, dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PrfItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PrfItemDao) Columns() PrfItemColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PrfItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PrfItemDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PrfItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
