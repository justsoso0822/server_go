package service

import (
	"context"
	"fmt"
	"time"

	"server_go/dao"
	"server_go/dao/model"
	"server_go/tools"
)

func InitTasks(ctx context.Context, uid int64) ([]map[string]interface{}, error) {
	confStr, err := dao.GetUserDataValue(ctx, uid, "task_conf")
	if err != nil {
		return nil, err
	}
	if confStr == "" {
		_ = dao.InsertUserData(ctx, &model.UserData{UID: int32(uid), Key: "task_conf", Value: "4"})
		confStr = "4"
	}

	serList := tools.PickNumbers(confStr)
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
	minId, maxId, err := dao.GetPrfTaskMinMax(ctx, ser)
	if err != nil {
		return nil, err
	}
	if minId == 0 {
		return nil, fmt.Errorf("用户%d的任务类型%d没有数据", uid, ser)
	}

	row, err := dao.GetUserTask(ctx, uid, minId, maxId)
	if err != nil {
		return nil, err
	}

	now := int(time.Now().Unix())
	var taskId int
	needClear := false

	if row == nil {
		taskId = minId
		_ = dao.InsertUserTask(ctx, &model.UserTask{UID: int32(uid), Taskid: int32(taskId), Addtm: int32(now)})
	} else if row.Done != 0 {
		taskId = int(row.Taskid)
		if taskId >= maxId {
			taskId, err = dao.GetLoopStartPrfTask(ctx, ser)
		} else {
			taskId, err = dao.GetNextPrfTask(ctx, ser, taskId)
		}
		if err != nil {
			return nil, err
		}
		needClear = true
		_ = dao.InsertUserTask(ctx, &model.UserTask{UID: int32(uid), Taskid: int32(taskId), Addtm: int32(now)})
	} else {
		taskId = int(row.Taskid)
	}

	if needClear {
		_ = dao.DeleteDoneUserTasks(ctx, uid, minId, maxId)
	}
	if taskId == 0 {
		return nil, nil
	}

	t, err := dao.GetPrfTaskById(ctx, taskId)
	if err != nil || t == nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":         t.ID,
		"ser":        t.Ser,
		"tid":        t.Tid,
		"npc":        t.Npc,
		"start_loop": t.StartLoop,
	}, nil
}
