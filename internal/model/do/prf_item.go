// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PrfItem is the golang structure of table prf_item for DAO operations like Where/Data.
type PrfItem struct {
	g.Meta     `orm:"table:prf_item, do:true"`
	Id         any //
	Ser        any // 道具系列
	Type       any // 道具类型  0 普通道具  1 工厂道具  2 转换器  3 点击使用的道具  4 拖拽到目标使用的道具
	Cost       any // 工厂类型道具产出消耗的体力
	Lv         any // 道具等级
	From       any // 来源id
	Name       any // 道具名称
	Tips       any // 道具介绍
	Star       any // 该道具对应的绿星数量
	Next       any // 合成道具id
	AutoDrop   any // 工厂是否自动产出道具
	NeedOpen   any // 是否需要等待开启
	CoolDrop   any // 一次CD掉落个数
	CoolTime   any // 冷却时间
	CoolMaxnum any // 冷却时间累积次数
	CoolMoney  any // 冷却加速花费钻石（最大值）
	Sell       any // 出售价格-金币
	Use        any // 可直接使用的道具，boxid
	Die        any // 特殊-掉落完后变成什么  <0 直接消失  ==0 进入cd  >0 变成指定道具, 如果道具不存在也消失
	Exp        any // 该等级道具合成下一级道具时是否会产出经验
	Gold       any // 副本合成掉落金币
	Rare       any // 0=普通，1=稀有（出售时跳出提醒框）
	Tids       any // 转换器内可投入物品的id
	Count      any // 转换器内对应投入物品的数量
}
