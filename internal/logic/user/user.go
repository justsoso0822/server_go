package user

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"server_go/internal/consts"
	"server_go/internal/dao"
	"server_go/internal/logic/gamelog"
	"server_go/internal/logic/lock"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sUser struct{}

func init() {
	service.RegisterUser(&sUser{})
}

func (s *sUser) Login(ctx context.Context, uid int64, loginKey, openid, platform, version string) (map[string]interface{}, error) {
	return Login(ctx, uid, loginKey, openid, platform, version)
}
func (s *sUser) UpdateDiamond(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error {
	return UpdateDiamond(ctx, uid, cnt, ret, reason)
}
func (s *sUser) UpdateGold(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error {
	return UpdateGold(ctx, uid, cnt, ret, reason)
}
func (s *sUser) UpdateTili(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error {
	return UpdateTili(ctx, uid, cnt, ret, reason)
}
func (s *sUser) UpdateExp(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error {
	return UpdateExp(ctx, uid, cnt, ret, reason)
}
func (s *sUser) UpdateStar(ctx context.Context, uid, cnt int64, ret map[string]interface{}, reason string) error {
	return UpdateStar(ctx, uid, cnt, ret, reason)
}
func (s *sUser) UpdateRes(ctx context.Context, uid int64, items interface{}, ret map[string]interface{}, reason string) error {
	return UpdateRes(ctx, uid, items, ret, reason)
}
func (s *sUser) UpdateItemsOther(ctx context.Context, uid int64, items map[int]int, ret map[string]interface{}, reason string) error {
	return UpdateItemsOther(ctx, uid, items, ret, reason)
}
func (s *sUser) GetUser(ctx context.Context, uid int64) (map[string]interface{}, error) {
	return GetUser(ctx, uid)
}
func (s *sUser) GetUserRes(ctx context.Context, uid int64) (map[string]interface{}, error) {
	return GetUserRes(ctx, uid)
}

// Login handles user login: find or create user, write login key, return full state.
func Login(ctx context.Context, uid int64, loginKey, openid, platform, version string) (g.Map, error) {
	if openid == "" {
		return g.Map{"code": -1, "msg": "参数错误"}, nil
	}

	ret := g.Map{"uid": uid}

	user, err := GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if g.NewVar(user["platform"]).String() != platform || g.NewVar(user["openid"]).String() != openid {
			return g.Map{"code": -1035, "msg": "账号信息不匹配"}, nil
		}
		ret["newbie"] = 0
		ret["user"] = user
	} else {
		ret["newbie"] = 1
		nowDay := gtime.Now().StartOfDay().Timestamp()
		err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, e := tx.Ctx(ctx).Exec(
				"INSERT INTO `user` (`uid`, `platform`, `openid`) VALUES (?, ?, ?)",
				uid, platform, openid,
			)
			if e != nil {
				return e
			}
			_, e = tx.Ctx(ctx).Exec(
				"INSERT INTO `user_res` (`uid`, `gold`, `diamond`, `star`, `tili`, `tili_time`, `exp`, `level`, `day_time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
				uid, 200, 100, 0, 100, 0, 0, 1, nowDay,
			)
			return e
		})
		if err != nil {
			return nil, err
		}
		ret["user"] = g.Map{
			"uid": uid, "platform": platform, "openid": openid,
			"sid": 0, "born": gtime.Now().Format("Y-m-d H:i:s"),
		}
	}

	// Log login (fire-and-forget)
	go func() {
		_, _ = g.DB().Exec(context.Background(),
			"INSERT INTO `_log_login` (`uid`, `platform`) VALUES (?, ?)", uid, platform)
	}()

	// Upsert login key
	_, err = g.DB().Exec(ctx,
		"INSERT INTO `user_loginkey` (`uid`, `key`, `ver`, `time`) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE `key` = VALUES(`key`), `ver` = VALUES(`ver`), `time` = VALUES(`time`)",
		uid, loginKey, version, gtime.Timestamp(),
	)
	if err != nil {
		return nil, err
	}

	// Fetch additional data
	datas, err := dao.UserData.Ctx(ctx).Where("uid", uid).All()
	if err != nil {
		return nil, err
	}
	ret["datas"] = datas

	gmVal, err := dao.SysGm.Ctx(ctx).Where("uid", uid).Value("uid")
	if err != nil {
		return nil, err
	}
	if gmVal.IsEmpty() {
		ret["gm"] = 0
	} else {
		ret["gm"] = 1
	}

	items, err := dao.UserItem.Ctx(ctx).Where("uid", uid).All()
	if err != nil {
		return nil, err
	}
	ret["items"] = items

	res, err := GetUserRes(ctx, uid)
	if err != nil {
		return nil, err
	}
	ret["res"] = res

	config, err := dao.MemConfig.Ctx(ctx).All()
	if err != nil {
		return nil, err
	}
	ret["config"] = config

	return ret, nil
}

// --- Resource update functions ---

func UpdateDiamond(ctx context.Context, uid, cnt int64, ret g.Map, reason string) error {
	lockKey := fmt.Sprintf("update_diamond:%d", uid)
	return updateUserResField(ctx, uid, "diamond", cnt, ret, reason, lockKey, "__add_diamond", "钻石")
}

func UpdateGold(ctx context.Context, uid, cnt int64, ret g.Map, reason string) error {
	lockKey := fmt.Sprintf("update_gold:%d", uid)
	return updateUserResField(ctx, uid, "gold", cnt, ret, reason, lockKey, "__add_gold", "金币")
}

func UpdateTili(ctx context.Context, uid, cnt int64, ret g.Map, reason string) error {
	lockKey := fmt.Sprintf("update_tili:%d", uid)
	return updateUserResField(ctx, uid, "tili", cnt, ret, reason, lockKey, "__add_tili", "体力")
}

func UpdateExp(ctx context.Context, uid, cnt int64, ret g.Map, reason string) error {
	lockKey := fmt.Sprintf("update_exp:%d", uid)
	return updateUserResField(ctx, uid, "exp", cnt, ret, reason, lockKey, "__add_exp", "经验")
}

func UpdateStar(ctx context.Context, uid, cnt int64, ret g.Map, reason string) error {
	lockKey := fmt.Sprintf("update_star:%d", uid)
	return updateUserResField(ctx, uid, "star", cnt, ret, reason, lockKey, "__add_star", "星星")
}

func updateUserResField(ctx context.Context, uid int64, field string, cnt int64, ret g.Map, reason, lockKey, addKey, resName string) error {
	token, err := lock.Lock(ctx, lockKey)
	if err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("系统繁忙，请稍后再试")
	}
	defer func() { _ = lock.Unlock(ctx, lockKey, token) }()

	res, err := GetUserRes(ctx, uid)
	if err != nil {
		return err
	}
	if res == nil {
		return fmt.Errorf("用户资源不存在")
	}

	oldCnt := g.NewVar(res[field]).Int64()
	newCnt := oldCnt + cnt
	if newCnt < 0 {
		newCnt = 0
	}
	if newCnt == oldCnt {
		return nil
	}

	_, err = dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{field: newCnt}).Update()
	if err != nil {
		gamelog.Log(ctx, uid, fmt.Sprintf("更新用户资源失败 %s %d %s %v", field, cnt, reason, err))
		return err
	}

	res[field] = newCnt
	ret["res"] = res

	prev := g.NewVar(ret[addKey]).Int64()
	ret[addKey] = (newCnt - oldCnt) + prev

	gamelog.TraceRes(ctx, uid, oldCnt, newCnt, resName, reason)
	return nil
}

// UpdateRes parses resource items and applies all changes.
func UpdateRes(ctx context.Context, uid int64, items interface{}, ret g.Map, reason string) error {
	resList := ParseRes(items)
	if len(resList) == 0 {
		return nil
	}

	var addDiamond, addGold, addTili, addExp, addStar int64
	addItemsOther := make(map[int]int)

	for _, r := range resList {
		switch r.Type {
		case consts.ResTypeDiamond:
			addDiamond += int64(r.Cnt)
		case consts.ResTypeGold:
			addGold += int64(r.Cnt)
		case consts.ResTypeTili:
			addTili += int64(r.Cnt)
		case consts.ResTypeExp:
			addExp += int64(r.Cnt)
		case consts.ResTypeStar:
			addStar += int64(r.Cnt)
		case consts.ResTypeItemOther:
			addItemsOther[r.Id] += r.Cnt
		default:
			return fmt.Errorf("不支持的资源类型 %d", r.Type)
		}
	}

	if addDiamond != 0 {
		if err := UpdateDiamond(ctx, uid, addDiamond, ret, reason); err != nil {
			return err
		}
	}
	if addGold != 0 {
		if err := UpdateGold(ctx, uid, addGold, ret, reason); err != nil {
			return err
		}
	}
	if addTili != 0 {
		if err := UpdateTili(ctx, uid, addTili, ret, reason); err != nil {
			return err
		}
	}
	if addExp != 0 {
		if err := UpdateExp(ctx, uid, addExp, ret, reason); err != nil {
			return err
		}
	}
	if addStar != 0 {
		if err := UpdateStar(ctx, uid, addStar, ret, reason); err != nil {
			return err
		}
	}
	if len(addItemsOther) > 0 {
		if err := UpdateItemsOther(ctx, uid, addItemsOther, ret, reason); err != nil {
			return err
		}
	}
	return nil
}

// UpdateItemsOther batch-updates non-grid items with locking and transaction.
func UpdateItemsOther(ctx context.Context, uid int64, items map[int]int, ret g.Map, reason string) error {
	if len(items) == 0 {
		return nil
	}
	lockKey := fmt.Sprintf("update_items_other:%d", uid)
	token, err := lock.Lock(ctx, lockKey)
	if err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("系统繁忙，请稍后再试")
	}
	defer func() { _ = lock.Unlock(ctx, lockKey, token) }()

	ids := make([]int, 0, len(items))
	for id := range items {
		if id != 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		rows, e := dao.UserItem.Ctx(ctx).TX(tx).Where("uid", uid).WhereIn("tid", ids).All()
		if e != nil {
			return e
		}
		oldCntMap := make(map[int]int)
		for _, row := range rows {
			oldCntMap[row["tid"].Int()] = row["cnt"].Int()
		}

		var valuesSql []string
		var upsertArgs []interface{}
		type updateInfo struct {
			Id, Delta, OldCnt, NewCnt int
		}
		var updates []updateInfo

		for _, id := range ids {
			delta := items[id]
			if delta == 0 {
				continue
			}
			oldCnt := oldCntMap[id]
			newCnt := oldCnt + delta
			if newCnt < 0 {
				newCnt = 0
			}
			valuesSql = append(valuesSql, "(?, ?, ?)")
			upsertArgs = append(upsertArgs, uid, id, newCnt)
			updates = append(updates, updateInfo{Id: id, Delta: delta, OldCnt: oldCnt, NewCnt: newCnt})
		}

		if len(valuesSql) > 0 {
			sql := fmt.Sprintf(
				"INSERT INTO `user_item` (`uid`, `tid`, `cnt`) VALUES %s ON DUPLICATE KEY UPDATE `cnt` = VALUES(`cnt`)",
				strings.Join(valuesSql, ","),
			)
			if _, e = tx.Ctx(ctx).Exec(sql, upsertArgs...); e != nil {
				return e
			}
		}

		retItems, _ := ret["items"].([]g.Map)
		if retItems == nil {
			retItems = []g.Map{}
		}
		addItemsOther, _ := ret["__add_items_other"].([]g.Map)
		if addItemsOther == nil {
			addItemsOther = []g.Map{}
		}

		for _, u := range updates {
			found := false
			for i, item := range retItems {
				if g.NewVar(item["tid"]).Int() == u.Id {
					retItems[i]["cnt"] = u.NewCnt
					found = true
					break
				}
			}
			if !found {
				retItems = append(retItems, g.Map{"uid": uid, "tid": u.Id, "cnt": u.NewCnt})
			}
			addItemsOther = append(addItemsOther, g.Map{"type": 8, "id": u.Id, "cnt": u.Delta})
			gamelog.TraceRes(ctx, uid, int64(u.OldCnt), int64(u.NewCnt), fmt.Sprintf("道具-%d", u.Id), reason)
		}
		ret["items"] = retItems
		ret["__add_items_other"] = addItemsOther
		return nil
	})

	if err != nil {
		gamelog.Log(ctx, uid, fmt.Sprintf("批量更新用户物品失败 %s %v", reason, err))
	}
	return err
}

// --- User/UserRes getters ---

func GetUser(ctx context.Context, uid int64) (g.Map, error) {
	row, err := dao.User.Ctx(ctx).Where("uid", uid).One()
	if err != nil {
		return nil, err
	}
	if row.IsEmpty() {
		return nil, nil
	}
	return row.Map(), nil
}

func GetUserRes(ctx context.Context, uid int64) (g.Map, error) {
	row, err := dao.UserRes.Ctx(ctx).Where("uid", uid).One()
	if err != nil {
		return nil, err
	}
	if row.IsEmpty() {
		return nil, nil
	}

	res := row.Map()
	nowDay := gtime.Now().StartOfDay().Timestamp()
	if g.NewVar(res["day_time"]).Int64() != nowDay {
		_, _ = dao.UserRes.Ctx(ctx).Where("uid", uid).Data(g.Map{"day_conf": "", "day_time": nowDay}).Update()
		res["day_conf"] = ""
		res["day_time"] = nowDay
	}
	res["now"] = gtime.TimestampMilli()
	return res, nil
}

// ParseRes parses resource items from string "type,id,cnt,..." or slice.
func ParseRes(items interface{}) []consts.ResItem {
	switch v := items.(type) {
	case []consts.ResItem:
		return v
	case string:
		return parseResString(v)
	default:
		return nil
	}
}

func parseResString(s string) []consts.ResItem {
	nums := pickNumbers(s)
	if len(nums) == 0 || len(nums)%3 != 0 {
		return nil
	}
	result := make([]consts.ResItem, 0, len(nums)/3)
	for i := 0; i < len(nums); i += 3 {
		result = append(result, consts.ResItem{Type: nums[i], Id: nums[i+1], Cnt: nums[i+2]})
	}
	return result
}

func pickNumbers(s string) []int {
	var result []int
	current := ""
	for _, c := range s {
		if (c >= '0' && c <= '9') || c == '-' || c == '+' || c == '.' {
			current += string(c)
		} else {
			if current != "" {
				if n, err := strconv.Atoi(current); err == nil {
					result = append(result, n)
				} else if f, err := strconv.ParseFloat(current, 64); err == nil {
					result = append(result, int(math.Floor(f)))
				}
				current = ""
			}
		}
	}
	if current != "" {
		if n, err := strconv.Atoi(current); err == nil {
			result = append(result, n)
		} else if f, err := strconv.ParseFloat(current, 64); err == nil {
			result = append(result, int(math.Floor(f)))
		}
	}
	return result
}