package task

import (
	"context"
	"fmt"

	"server_go/internal/dao"
	"server_go/internal/logic/user"
	"server_go/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sTask struct{}

func init() {
	service.RegisterTask(&sTask{})
}

func (s *sTask) InitTasks(ctx context.Context, uid int64) ([]map[string]interface{}, error) {
	taskConf, err := dao.UserData.Ctx(ctx).Where("uid", uid).Where("key", "task_conf").Value("value")
	if err != nil {
		return nil, err
	}

	confStr := taskConf.String()
	if confStr == "" {
		_, _ = dao.UserData.Ctx(ctx).Data(g.Map{"uid": uid, "key": "task_conf", "value": "4"}).Insert()
		confStr = "4"
	}

	serList := user.PickNumbers(confStr)
	var arr []map[string]interface{}

	for _, ser := range serList {
		task, e := getOneTask(ctx, uid, ser)
		if e != nil {
			return nil, e
		}
		if task != nil {
			task["uid"] = uid
			arr = append(arr, task)
		}
	}

	return arr, nil
}

func getOneTask(ctx context.Context, uid int64, ser int) (map[string]interface{}, error) {
	minMax, err := getTaskSerMinMax(ctx, ser)
	if err != nil {
		return nil, err
	}
	if minMax == nil {
		return nil, fmt.Errorf("用户%d的任务类型%d没有数据", uid, ser)
	}
	minId, maxId := minMax[0], minMax[1]

	row, err := dao.UserTask.Ctx(ctx).Where("uid", uid).
		WhereGTE("taskid", minId).WhereLTE("taskid", maxId).
		Limit(1).One()
	if err != nil {
		return nil, err
	}

	var taskId int
	needClear := false
	nowSec := gtime.Timestamp()

	if row.IsEmpty() {
		taskId = minId
		_, err = dao.UserTask.Ctx(ctx).Data(g.Map{
			"uid": uid, "taskid": taskId, "addtm": nowSec, "done": 0, "donetm": 0,
		}).Insert()
		if err != nil {
			return nil, err
		}
	} else if row["done"].Int() != 0 {
		taskId = row["taskid"].Int()
		if taskId >= maxId {
			v, e := dao.PrfTask.Ctx(ctx).Where("ser", ser).Where("start_loop", 1).
				OrderAsc("id").Limit(1).Value("id")
			if e != nil {
				return nil, e
			}
			taskId = v.Int()
		} else {
			v, e := dao.PrfTask.Ctx(ctx).Where("ser", ser).WhereGT("id", taskId).
				OrderAsc("id").Limit(1).Value("id")
			if e != nil {
				return nil, e
			}
			taskId = v.Int()
		}
		needClear = true
		_, err = dao.UserTask.Ctx(ctx).Data(g.Map{
			"uid": uid, "taskid": taskId, "addtm": nowSec, "done": 0, "donetm": 0,
		}).Insert()
		if err != nil {
			return nil, err
		}
	} else {
		taskId = row["taskid"].Int()
	}

	if needClear {
		_, _ = dao.UserTask.Ctx(ctx).Where("uid", uid).
			WhereGTE("taskid", minId).WhereLTE("taskid", maxId).
			Where("done", 1).Delete()
	}

	if taskId == 0 {
		return nil, nil
	}

	taskRow, err := dao.PrfTask.Ctx(ctx).Where("id", taskId).One()
	if err != nil {
		return nil, err
	}
	if taskRow.IsEmpty() {
		return nil, nil
	}
	return taskRow.Map(), nil
}

func getTaskSerMinMax(ctx context.Context, ser int) ([]int, error) {
	minVal, err := dao.PrfTask.Ctx(ctx).Where("ser", ser).Min("id")
	if err != nil {
		return nil, err
	}
	if minVal == 0 {
		return nil, nil
	}
	maxVal, err := dao.PrfTask.Ctx(ctx).Where("ser", ser).Max("id")
	if err != nil {
		return nil, err
	}
	if maxVal == 0 {
		return nil, nil
	}
	return []int{int(minVal), int(maxVal)}, nil
}
