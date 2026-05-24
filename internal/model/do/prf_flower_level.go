// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PrfFlowerLevel is the golang structure of table prf_flower_level for DAO operations like Where/Data.
type PrfFlowerLevel struct {
	g.Meta       `orm:"table:prf_flower_level, do:true"`
	Id           any //
	Flower       any // 花id
	Level        any // 花-等级
	SeedUp       any // 升到下级需要的种子，最大等级填-1
	CoinUp       any // 升到下级需要的金币，最大等级填-1
	ReapExp      any // 收割经验
	ReapRound    any // 收获次数
	ReapInterval any // 收获间隔
	SeedDrop     any // 种子掉落数量
	SeedDropRate any // 种子掉落概率，填100代表1%概率
	ReapCnt      any // 一次收获几个
	FlowerMin    any // 摘花保底数量
}
