package api

import (
	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/server/core"
	"gopkg.in/mgo.v2/bson"
)

// TaskController Task相关
type TaskController struct {
	atopAPIController
}

// NewTask 创建新的Task任务。
// Task信息需要正常，如何新建任务成功将返回任务ID，否则返回错误信息。
//
//@router /task/new [post]
func (t *TaskController) NewTask() {
	task := &models.Task{}
	if err := t.UnmarshalBody(task); err != nil {
		t.OutputError(err)
		return
	}
	if err := core.TaskMgt.NewTask(task); err != nil {
		t.OutputError(err)
		return
	}
	t.OutputSuccess(map[string]string{"taskID": task.ID.Hex()})
}

// Start 执行指定任务，并返回任务推送结果。
// 启动任务前需要保证任务信息存在，且执行命令的服务器Agent也在线。
//
// @router /task/start [post]
func (t *TaskController) Start() {
	data := make(map[string]string)
	if err := t.UnmarshalBody(&data); err != nil {
		t.OutputError(err)
		return
	}
	if taskID, ok := data["taskID"]; !ok {
		t.OutputError("缺少参数taskID")
	} else if taskID == "" {
		t.OutputError("taskID不能为空")
	} else {
		err := core.TaskMgt.StartTask(bson.ObjectIdHex(taskID))
		t.OutputDoResult(err)
	}
}
