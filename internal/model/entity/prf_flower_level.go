// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PrfFlowerLevel is the golang structure for table prf_flower_level.
type PrfFlowerLevel struct {
	Id           int    `json:"id"           orm:"id"             description:""`
	Flower       int    `json:"flower"       orm:"flower"         description:"花id"`
	Level        int    `json:"level"        orm:"level"          description:"花-等级"`
	SeedUp       int    `json:"seedUp"       orm:"seed_up"        description:"升到下级需要的种子，最大等级填-1"`
	CoinUp       int    `json:"coinUp"       orm:"coin_up"        description:"升到下级需要的金币，最大等级填-1"`
	ReapExp      int    `json:"reapExp"      orm:"reap_exp"       description:"收割经验"`
	ReapRound    int    `json:"reapRound"    orm:"reap_round"     description:"收获次数"`
	ReapInterval int    `json:"reapInterval" orm:"reap_interval"  description:"收获间隔"`
	SeedDrop     int    `json:"seedDrop"     orm:"seed_drop"      description:"种子掉落数量"`
	SeedDropRate string `json:"seedDropRate" orm:"seed_drop_rate" description:"种子掉落概率，填100代表1%概率"`
	ReapCnt      int    `json:"reapCnt"      orm:"reap_cnt"       description:"一次收获几个"`
	FlowerMin    int    `json:"flowerMin"    orm:"flower_min"     description:"摘花保底数量"`
}
