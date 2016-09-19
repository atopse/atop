package task

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/db"
	"github.com/ysqi/atop/common/models"
)

// UpdateTaskStatus 更新任务状态
func UpdateTaskStatus(taskID bson.ObjectId, newStatus models.TaskStatus) error {
	if taskID == "" {
		return common.ErrBadBody{Message: "缺少参数: taskId"}
	} else if newStatus == "" {
		return common.ErrBadBody{Message: "缺少参数: status"}
	}

	m := bson.M{"status": newStatus}
	switch newStatus {
	case models.TaskStatusStarted:
		m["options.started"] = time.Now()
	case models.TaskStatusCompleted:
		m["options.completed"] = time.Now()
	case models.TaskStatusErrorDown:
		m["options.completed"] = time.Now()
	case models.TaskStatusProcessing:
		//进行中无需登记
	default:
		return common.ErrBadBody{Message: "该TaskStatus=%s不支持更新"}
	}

	return db.Do(func(dataBase *mgo.Database) error {
		var oldStatus models.TaskStatus
		err := dataBase.C("task").FindId(taskID).Select(bson.M{"status": 1}).One(&oldStatus)
		if err == mgo.ErrNotFound {
			return common.ErrBadBody{Message: err.Error()}
		}
		if newStatus == oldStatus {
			return nil
		}
		if oldStatus == models.TaskStatusCompleted || oldStatus == models.TaskStatusErrorDown {
			return nil
		}

		return dataBase.C("task").UpdateId(taskID, bson.M{
			"$set": m,
		})
	})
}

// NewProcess 更新任务处理进度
func NewProcess(msg *models.Msg) error {
	process, ok := msg.Content.(*models.CmdExecProcess)
	if !ok {
		return common.BadBodyErrf("非任务处理进度消息，msg.Content不是CmdExecProcess,而是%v", msg.Content)
	}
	if process == nil {
		return common.BadBodyErr("信息不能为空")
	}
	if process.CommandID == "" {
		return common.BadBodyErr("缺失TaskID(CommandID)")
	}

	if process.Tag == "" {
		return common.BadBodyErr("任务进行消息Tag不能为空")
	}

	taskID := process.CommandID
	if db.RecordIsExistByID("task", taskID) == false {
		return common.BadBodyErr("指定的Task不存在")
	}

	var newStatus models.TaskStatus
	if process.Tag == "newStatus" {
		s := process.Body.(string)
		if s == "" {
			return common.BadBodyErr("Tag=newStatus时，新状态内容不能为空字符串")
		}
		newStatus = models.TaskStatus(s)
		if newStatus != models.TaskStatusErrorDown && newStatus != models.TaskStatusProcessing && newStatus != models.TaskStatusCompleted {
			return fmt.Errorf("无法处理的任务状态%q", newStatus)
		}
	} else if process.Tag == "error" {
		newStatus = models.TaskStatusErrorDown
	}

	//started,processing,stopped,completed
	//先更新状态，再记录日志
	if err := UpdateTaskStatus(taskID, newStatus); err != nil {
		return err
	}
	taskLog := &models.TaskLog{
		TaskID:         process.CommandID,
		OccurrenceTime: msg.Created,
		Content:        process.Body,
	}
	handle.PushTaskLog(taskLog)
	return nil
}

// NewTaskLog 记录任务日志
func NewTaskLog(tlog models.TaskLog) error {
	if tlog.TaskID == "" {
		return common.BadBodyErr("任务日志记录TaskID不能为空")
	}
	if tlog.Content == nil {
		return common.BadBodyErr("任务日志内容不能为空")
	}
	return db.Do(func(dataBase *mgo.Database) error {
		if c, err := dataBase.C("task").FindId(tlog.TaskID).Count(); err != nil {
			return err
		} else if c == 0 {
			return common.BadBodyErrf("任务日志记录Task<%s>不存在,不允许添加任务日志", tlog.TaskID)
		}
		return dataBase.C("tasklog").Insert(tlog)
	})
}
