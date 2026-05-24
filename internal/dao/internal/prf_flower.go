// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"server_go/internal/autodb"

	"github.com/gogf/gf/v2/database/gdb"
)

// PrfFlowerDao is the data access object for the table prf_flower.
type PrfFlowerDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PrfFlowerColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PrfFlowerColumns defines and stores column names for the table prf_flower.
type PrfFlowerColumns struct {
	Id                   string //
	Name                 string // 花的名字
	Des                  string // 花语
	From                 string // 0=升级解锁 1=测试用 2=活动获得 3=鲜花礼包  4=徐霞客处获得  5=VIP专享  6=公会获得  7=敬请期待 8=已绝版
	FromTip              string // 来源提示
	Index                string // 排序
	SellRes              string // 出售获得的资源，类型,id,数量
	Water                string // 种植消耗的水滴
	Cost                 string // 培育所需资源
	Cd                   string // 培育所需时间
	Qua                  string // 品质。1=绿色，2=蓝色，3=紫色，4=红色，5=金色
	Charm                string // 魅力值
	Pic0                 string // 培育房-种子、收割-掉落的种子、鲜花升级-卡片种子
	Pic1                 string // 种地-刚刚播种
	Pic2                 string // 种地-长了一半
	Pic3                 string // 种地-完全长成、徐霞客、社团种植花盆内
	Pic4                 string // 收割-掉在地上
	Pic5                 string // 培育房-培育花台缩略图、仓库-缩略图、播种选择页、鲜花订单、插花、社团种植收取、好友交易、按住播种
	Pic6                 string // 花谱-外层显示、花市订单
	Pic7                 string // 花谱-带花瓶
	PicLand              string // 土地图片
	SpineFlowerLand      string // 种在地里的花动画
	SpineLand            string // 地块动画
	SpineFlowerBook      string // 图鉴花动画
	SpineFlowerBookExtra string // 图鉴花动画-附加
	NoBottle             string // 不显示瓶子
}

// prfFlowerColumns holds the columns for the table prf_flower.
var prfFlowerColumns = PrfFlowerColumns{
	Id:                   "id",
	Name:                 "name",
	Des:                  "des",
	From:                 "from",
	FromTip:              "from_tip",
	Index:                "index",
	SellRes:              "sell_res",
	Water:                "water",
	Cost:                 "cost",
	Cd:                   "cd",
	Qua:                  "qua",
	Charm:                "charm",
	Pic0:                 "pic_0",
	Pic1:                 "pic_1",
	Pic2:                 "pic_2",
	Pic3:                 "pic_3",
	Pic4:                 "pic_4",
	Pic5:                 "pic_5",
	Pic6:                 "pic_6",
	Pic7:                 "pic_7",
	PicLand:              "pic_land",
	SpineFlowerLand:      "spine_flower_land",
	SpineLand:            "spine_land",
	SpineFlowerBook:      "spine_flower_book",
	SpineFlowerBookExtra: "spine_flower_book_extra",
	NoBottle:             "no_bottle",
}

// NewPrfFlowerDao creates and returns a new DAO object for table data access.
func NewPrfFlowerDao(handlers ...gdb.ModelHandler) *PrfFlowerDao {
	return &PrfFlowerDao{
		group:    "default",
		table:    "prf_flower",
		columns:  prfFlowerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the raw database management object of the current DAO using request context.
func (dao *PrfFlowerDao) DB(ctx context.Context) gdb.DB {
	return autodb.DB(ctx, dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PrfFlowerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PrfFlowerDao) Columns() PrfFlowerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PrfFlowerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PrfFlowerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PrfFlowerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
