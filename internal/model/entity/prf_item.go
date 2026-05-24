// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PrfItem is the golang structure for table prf_item.
type PrfItem struct {
	Id         int    `json:"id"         orm:"id"          description:""`
	Ser        int    `json:"ser"        orm:"ser"         description:"道具系列"`
	Type       int    `json:"type"       orm:"type"        description:"道具类型  0 普通道具  1 工厂道具  2 转换器  3 点击使用的道具  4 拖拽到目标使用的道具"`
	Cost       int    `json:"cost"       orm:"cost"        description:"工厂类型道具产出消耗的体力"`
	Lv         int    `json:"lv"         orm:"lv"          description:"道具等级"`
	From       string `json:"from"       orm:"from"        description:"来源id"`
	Name       string `json:"name"       orm:"name"        description:"道具名称"`
	Tips       string `json:"tips"       orm:"tips"        description:"道具介绍"`
	Star       int    `json:"star"       orm:"star"        description:"该道具对应的绿星数量"`
	Next       int    `json:"next"       orm:"next"        description:"合成道具id"`
	AutoDrop   int    `json:"autoDrop"   orm:"auto_drop"   description:"工厂是否自动产出道具"`
	NeedOpen   int    `json:"needOpen"   orm:"need_open"   description:"是否需要等待开启"`
	CoolDrop   int    `json:"coolDrop"   orm:"cool_drop"   description:"一次CD掉落个数"`
	CoolTime   int    `json:"coolTime"   orm:"cool_time"   description:"冷却时间"`
	CoolMaxnum int    `json:"coolMaxnum" orm:"cool_maxnum" description:"冷却时间累积次数"`
	CoolMoney  int    `json:"coolMoney"  orm:"cool_money"  description:"冷却加速花费钻石（最大值）"`
	Sell       int    `json:"sell"       orm:"sell"        description:"出售价格-金币"`
	Use        int    `json:"use"        orm:"use"         description:"可直接使用的道具，boxid"`
	Die        int    `json:"die"        orm:"die"         description:"特殊-掉落完后变成什么  <0 直接消失  ==0 进入cd  >0 变成指定道具, 如果道具不存在也消失"`
	Exp        int    `json:"exp"        orm:"exp"         description:"该等级道具合成下一级道具时是否会产出经验"`
	Gold       int    `json:"gold"       orm:"gold"        description:"副本合成掉落金币"`
	Rare       int    `json:"rare"       orm:"rare"        description:"0=普通，1=稀有（出售时跳出提醒框）"`
	Tids       string `json:"tids"       orm:"tids"        description:"转换器内可投入物品的id"`
	Count      string `json:"count"      orm:"count"       description:"转换器内对应投入物品的数量"`
}
