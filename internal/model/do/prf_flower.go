// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PrfFlower is the golang structure of table prf_flower for DAO operations like Where/Data.
type PrfFlower struct {
	g.Meta               `orm:"table:prf_flower, do:true"`
	Id                   any //
	Name                 any // 花的名字
	Des                  any // 花语
	From                 any // 0=升级解锁 1=测试用 2=活动获得 3=鲜花礼包  4=徐霞客处获得  5=VIP专享  6=公会获得  7=敬请期待 8=已绝版
	FromTip              any // 来源提示
	Index                any // 排序
	SellRes              any // 出售获得的资源，类型,id,数量
	Water                any // 种植消耗的水滴
	Cost                 any // 培育所需资源
	Cd                   any // 培育所需时间
	Qua                  any // 品质。1=绿色，2=蓝色，3=紫色，4=红色，5=金色
	Charm                any // 魅力值
	Pic0                 any // 培育房-种子、收割-掉落的种子、鲜花升级-卡片种子
	Pic1                 any // 种地-刚刚播种
	Pic2                 any // 种地-长了一半
	Pic3                 any // 种地-完全长成、徐霞客、社团种植花盆内
	Pic4                 any // 收割-掉在地上
	Pic5                 any // 培育房-培育花台缩略图、仓库-缩略图、播种选择页、鲜花订单、插花、社团种植收取、好友交易、按住播种
	Pic6                 any // 花谱-外层显示、花市订单
	Pic7                 any // 花谱-带花瓶
	PicLand              any // 土地图片
	SpineFlowerLand      any // 种在地里的花动画
	SpineLand            any // 地块动画
	SpineFlowerBook      any // 图鉴花动画
	SpineFlowerBookExtra any // 图鉴花动画-附加
	NoBottle             any // 不显示瓶子
}
