// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"server_go/internal/autodb"

	"github.com/gogf/gf/v2/database/gdb"
)

// PrfFlowerLevelDao is the data access object for the table prf_flower_level.
type PrfFlowerLevelDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  PrfFlowerLevelColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// PrfFlowerLevelColumns defines and stores column names for the table prf_flower_level.
type PrfFlowerLevelColumns struct {
	Id           string //
	Flower       string // 花id
	Level        string // 花-等级
	SeedUp       string // 升到下级需要的种子，最大等级填-1
	CoinUp       string // 升到下级需要的金币，最大等级填-1
	ReapExp      string // 收割经验
	ReapRound    string // 收获次数
	ReapInterval string // 收获间隔
	SeedDrop     string // 种子掉落数量
	SeedDropRate string // 种子掉落概率，填100代表1%概率
	ReapCnt      string // 一次收获几个
	FlowerMin    string // 摘花保底数量
}

// prfFlowerLevelColumns holds the columns for the table prf_flower_level.
var prfFlowerLevelColumns = PrfFlowerLevelColumns{
	Id:           "id",
	Flower:       "flower",
	Level:        "level",
	SeedUp:       "seed_up",
	CoinUp:       "coin_up",
	ReapExp:      "reap_exp",
	ReapRound:    "reap_round",
	ReapInterval: "reap_interval",
	SeedDrop:     "seed_drop",
	SeedDropRate: "seed_drop_rate",
	ReapCnt:      "reap_cnt",
	FlowerMin:    "flower_min",
}

// NewPrfFlowerLevelDao creates and returns a new DAO object for table data access.
func NewPrfFlowerLevelDao(handlers ...gdb.ModelHandler) *PrfFlowerLevelDao {
	return &PrfFlowerLevelDao{
		group:    "default",
		table:    "prf_flower_level",
		columns:  prfFlowerLevelColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the raw database management object of the current DAO using request context.
func (dao *PrfFlowerLevelDao) DB(ctx context.Context) gdb.DB {
	return autodb.DB(ctx, dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PrfFlowerLevelDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PrfFlowerLevelDao) Columns() PrfFlowerLevelColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PrfFlowerLevelDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PrfFlowerLevelDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PrfFlowerLevelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
