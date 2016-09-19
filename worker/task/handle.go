package task

import (
	"github.com/nsqio/go-nsq"
	"github.com/pquerna/ffjson/ffjson"
	"gopkg.in/mgo.v2/bson"

	"fmt"

	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/worker/app"
)

var handle = &MessageHandle{}

// MessageHandle NSQ消息处理
type MessageHandle struct {
}

// HandleNewStatus 更新任务状态
func (m *MessageHandle) HandleNewStatus(message *nsq.Message) error {
	body := struct {
		TaskID bson.ObjectId     `json:"taskId"`
		Status models.TaskStatus `json:"status"`
	}{}
	if err := ffjson.Unmarshal(message.Body, &body); err != nil {
		return common.ErrBadBody{Message: err.Error()}
	}
	return UpdateTaskStatus(body.TaskID, body.Status)
}

// HandelNewProcess 更新任务处理进度.
func (m *MessageHandle) HandelNewProcess(message *nsq.Message) error {
	msg := &models.Msg{}
	if err := ffjson.Unmarshal(message.Body, msg); err != nil {
		return common.ErrBadBody{Message: err.Error()}
	}
	err := NewProcess(msg)
	if err != nil {
		//如果处理出错，则需要将错误信息记录到任务日志中
		process, ok := msg.Content.(*models.CmdExecProcess)
		if !ok || process == nil {
			return err
		}
		taskLog := &models.TaskLog{
			TaskID:         process.CommandID,
			OccurrenceTime: msg.Created,
			Content: map[string]interface{}{
				"error":   fmt.Sprintf("任务进度处理失败，%s", err),
				"message": process.Body},
		}
		m.PushTaskLog(taskLog)
	}
	return err
}

// HandelNewTaskLog 处理任务日志
func (m *MessageHandle) HandelNewTaskLog(message *nsq.Message) error {
	tlog := models.TaskLog{}
	if err := ffjson.Unmarshal(message.Body, &tlog); err != nil {
		return common.ErrBadBody{Message: err.Error()}
	}
	return NewTaskLog(tlog)
}

// PushTaskLog 推送任务日志消息.
func (m *MessageHandle) PushTaskLog(tlog *models.TaskLog) {
	app.ProduceWorker.Publish("task.log.new", tlog)
}

func init() {
	app.AddHandle("task.status.new", handle.HandleNewStatus)
	app.AddHandle("task.log.new", handle.HandelNewProcess)
}
