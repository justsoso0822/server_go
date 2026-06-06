package service

import (
	"context"
	"fmt"
	"time"

	"server_go/dao"
	"server_go/dao/model"
	"server_go/tools"
)

func InitTasks(ctx context.Context, uid int64) ([]map[string]any, error) {
	confStr, err := dao.GetUserDataValue(ctx, uid, "task_conf")
	if err != nil {
		return nil, err
	}
	if confStr == "" {
		_ = dao.InsertUserData(ctx, &model.UserData{UID: uid, Key: "task_conf", Value: "4"})
		confStr = "4"
	}

	serList := tools.PickNumbers(confStr)
	var arr []map[string]any
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

func getOneTask(ctx context.Context, uid int64, ser int) (map[string]any, error) {
	minID, maxID, err := dao.GetPrfTaskMinMax(ctx, ser)
	if err != nil {
		return nil, err
	}
	if minID == 0 {
		return nil, fmt.Errorf("用户%d的任务类型%d没有数据", uid, ser)
	}

	row, err := dao.GetUserTask(ctx, uid, minID, maxID)
	if err != nil {
		return nil, err
	}

	now := int(time.Now().Unix())
	var taskID int
	needClear := false

	if row == nil {
		taskID = minID
		_ = dao.InsertUserTask(ctx, &model.UserTask{UID: uid, Taskid: int32(taskID), Addtm: int32(now)})
	} else if row.Done != 0 {
		taskID = int(row.Taskid)
		if taskID >= maxID {
			taskID, err = dao.GetLoopStartPrfTask(ctx, ser)
		} else {
			taskID, err = dao.GetNextPrfTask(ctx, ser, taskID)
		}
		if err != nil {
			return nil, err
		}
		needClear = true
		_ = dao.InsertUserTask(ctx, &model.UserTask{UID: uid, Taskid: int32(taskID), Addtm: int32(now)})
	} else {
		taskID = int(row.Taskid)
	}

	if needClear {
		_ = dao.DeleteDoneUserTasks(ctx, uid, minID, maxID)
	}
	if taskID == 0 {
		return nil, nil
	}

	t, err := dao.GetPrfTaskByID(ctx, taskID)
	if err != nil || t == nil {
		return nil, err
	}
	return map[string]any{
		"id":         t.ID,
		"ser":        t.Ser,
		"tid":        t.Tid,
		"npc":        t.Npc,
		"start_loop": t.StartLoop,
	}, nil
}
