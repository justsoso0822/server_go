// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PrfFlower is the golang structure for table prf_flower.
type PrfFlower struct {
	Id                   int    `json:"id"                   orm:"id"                      description:""`
	Name                 string `json:"name"                 orm:"name"                    description:"花的名字"`
	Des                  string `json:"des"                  orm:"des"                     description:"花语"`
	From                 int    `json:"from"                 orm:"from"                    description:"0=升级解锁 1=测试用 2=活动获得 3=鲜花礼包  4=徐霞客处获得  5=VIP专享  6=公会获得  7=敬请期待 8=已绝版"`
	FromTip              string `json:"fromTip"              orm:"from_tip"                description:"来源提示"`
	Index                int    `json:"index"                orm:"index"                   description:"排序"`
	SellRes              string `json:"sellRes"              orm:"sell_res"                description:"出售获得的资源，类型,id,数量"`
	Water                int    `json:"water"                orm:"water"                   description:"种植消耗的水滴"`
	Cost                 string `json:"cost"                 orm:"cost"                    description:"培育所需资源"`
	Cd                   int    `json:"cd"                   orm:"cd"                      description:"培育所需时间"`
	Qua                  int    `json:"qua"                  orm:"qua"                     description:"品质。1=绿色，2=蓝色，3=紫色，4=红色，5=金色"`
	Charm                int    `json:"charm"                orm:"charm"                   description:"魅力值"`
	Pic0                 string `json:"pic0"                 orm:"pic_0"                   description:"培育房-种子、收割-掉落的种子、鲜花升级-卡片种子"`
	Pic1                 string `json:"pic1"                 orm:"pic_1"                   description:"种地-刚刚播种"`
	Pic2                 string `json:"pic2"                 orm:"pic_2"                   description:"种地-长了一半"`
	Pic3                 string `json:"pic3"                 orm:"pic_3"                   description:"种地-完全长成、徐霞客、社团种植花盆内"`
	Pic4                 string `json:"pic4"                 orm:"pic_4"                   description:"收割-掉在地上"`
	Pic5                 string `json:"pic5"                 orm:"pic_5"                   description:"培育房-培育花台缩略图、仓库-缩略图、播种选择页、鲜花订单、插花、社团种植收取、好友交易、按住播种"`
	Pic6                 string `json:"pic6"                 orm:"pic_6"                   description:"花谱-外层显示、花市订单"`
	Pic7                 string `json:"pic7"                 orm:"pic_7"                   description:"花谱-带花瓶"`
	PicLand              string `json:"picLand"              orm:"pic_land"                description:"土地图片"`
	SpineFlowerLand      string `json:"spineFlowerLand"      orm:"spine_flower_land"       description:"种在地里的花动画"`
	SpineLand            string `json:"spineLand"            orm:"spine_land"              description:"地块动画"`
	SpineFlowerBook      string `json:"spineFlowerBook"      orm:"spine_flower_book"       description:"图鉴花动画"`
	SpineFlowerBookExtra string `json:"spineFlowerBookExtra" orm:"spine_flower_book_extra" description:"图鉴花动画-附加"`
	NoBottle             int    `json:"noBottle"             orm:"no_bottle"               description:"不显示瓶子"`
}
