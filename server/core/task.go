package core

import (
	"errors"
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	"github.com/otium/queue"

	"github.com/ysqi/atop/common/db"
	"github.com/ysqi/atop/common/models"
)

// TaskMgt 任务工作服务
var TaskMgt = &TaskService{}

// TaskService 任务工作服务
type TaskService struct {
	queue *queue.Queue
}

// GetTaskInfo 获取Task信息
func (t *TaskService) GetTaskInfo(taskID interface{}) (*models.Task, error) {
	task := &models.Task{}
	return task, db.Do(func(dataBase *mgo.Database) error {
		return dataBase.C("task").FindId(taskID).One(task)
	})
}

// NewTask 返回新建任务处理结果。
// 新建任务时，将再Command信息中记录TaskID
func (t *TaskService) NewTask(task *models.Task) error {
	if err := task.Verify(); err != nil {
		return err
	}
	task.ID = bson.NewObjectId()
	task.Cmd.ID = task.ID //保存任务ID
	if task.Options == nil {
		task.Options = make(map[string]interface{})
	}
	task.Status = models.TaskStatusNew
	return db.Insert("task", task)
}

// StartTask 开启任务
func (t *TaskService) StartTask(taskID bson.ObjectId) error {
	if taskID == "" {
		return errors.New("任务ID")
	}
	task, err := t.GetTaskInfo(taskID)
	if err != nil {
		return err
	}
	if task.Status != models.TaskStatusNew {
		return fmt.Errorf("任务状态<%s>为非初始状态,不允许重新允许", task.Status)
	}
	hasError := func(err error) error {
		t.PushLog(task.ID, err.Error())
		if err2 := t.UpdateTaskStatus(task.ID, models.TaskStatusErrorDown); err2 != nil {
			t.PushLog(task.ID, fmt.Sprintf("更新Task状态失败,%s", err))
		}
		return err
	}
	//1.获取Agent
	agent := AgentMgt.GetOnlineAgent(task.TargetIP, task.TargetIP2)
	if agent == nil {
		return hasError(fmt.Errorf("任务执行目标服务器[%s,%s]不存在", task.TargetIP, task.TargetIP2))
	}

	_, err = AgentMgt.HTTPDoRequest(agent, "post", "/command/exec", task.Cmd)
	if err != nil {
		return hasError(fmt.Errorf("任务推送给Agent:%s失败,%s", agent, err))
	}
	t.PushLog(task.ID, fmt.Sprintf("Agent:%s已受理任务", agent))
	if err := t.UpdateTaskStatus(task.ID, models.TaskStatusStarted); err != nil {
		t.PushLog(task.ID, fmt.Sprintf("更新Task状态失败,%s", err))
	}
	return nil
}

// UpdateTaskStatus 更新任务状态
func (t *TaskService) UpdateTaskStatus(taskID bson.ObjectId, to models.TaskStatus) error {
	if taskID == "" {
		return errors.New("入参taskID未空")
	}
	m := bson.M{"status": to}
	switch to {
	case models.TaskStatusStarted:
		m["options.started"] = time.Now()
	case models.TaskStatusCompleted:
		m["options.completed"] = time.Now()
	case models.TaskStatusErrorDown:
		m["options.completed"] = time.Now()
	case models.TaskStatusProcessing:
		//进行中无需登记
	default:
		return fmt.Errorf("该TaskStatus=%s不支持更新", to)
	}
	return db.Do(func(dataBase *mgo.Database) error {
		return dataBase.C("task").UpdateId(taskID, bson.M{
			"$set": m,
		})
	})
}

// PushTaskProcess 推送任务处理进度信息
func (t *TaskService) PushTaskProcess(msg *models.Msg) error {
	process, ok := msg.Content.(*models.CmdExecProcess)
	if !ok {
		return errors.New("非任务处理进度消息，msg.Content不是CmdExecProcess,而是：" + fmt.Sprintf("%v", msg.Content))
	}
	if process.CommandID == "" {
		return errors.New("缺失TaskID(CommandID)")
	}

	if process.Tag == "" {
		return errors.New("任务进行消息Tag不能为空")
	}

	var newStatus models.TaskStatus
	if process.Tag == "newStatus" {
		s := process.Body.(string)
		if s == "" {
			return fmt.Errorf("无法处理的任务状态%q", s)
		}
		newStatus = models.TaskStatus(s)
		if newStatus != models.TaskStatusErrorDown && newStatus != models.TaskStatusProcessing && newStatus != models.TaskStatusCompleted {
			return fmt.Errorf("无法处理的任务状态%q", newStatus)
		}
	} else if process.Tag == "error" {
		newStatus = models.TaskStatusErrorDown
	}

	taskID := process.CommandID
	if count, err := db.QueryCount("task", bson.M{"_id": taskID}); err != nil {
		return err
	} else if count == 0 {
		return errors.New("指定的Task不存在")
	}

	if newStatus != "" {
		//started,processing,stopped,completed
		//先更新状态，再记录日志
		err := t.UpdateTaskStatus(taskID, newStatus)
		if err != nil {
			taskLog := &models.TaskLog{
				TaskID:         process.CommandID,
				OccurrenceTime: msg.Created,
				Content:        map[string]interface{}{"error": err, "message": process.Body},
			}
			t.pushLog(taskLog)
			// t.PushLog(taskID, fmt.Sprintf("任务新状态:%s-%v,更新失败:%s", process.Tag, process.Body, err.Error()))
			return err
		}
	}
	taskLog := &models.TaskLog{
		TaskID:         process.CommandID,
		OccurrenceTime: msg.Created,
		Content:        process,
	}
	t.pushLog(taskLog)
	return nil
}

// PushLog 推送日志
func (t *TaskService) PushLog(taskID bson.ObjectId, content interface{}) {
	if content == "" {
		return
	}
	if taskID == "" {
		return
	}
	taskLog := &models.TaskLog{
		TaskID:         taskID,
		OccurrenceTime: time.Now(),
		Content:        content,
	}

	t.pushLog(taskLog)
}
func (t *TaskService) pushLog(log *models.TaskLog) {
	t.queue.Push(log)
}

// saveLog 保存task执行日志到DB.
// 收到log后，将日志保存到DB，如果保存失败，将消息继续PUSH到队列，重新处理.
func (t *TaskService) saveLog(log interface{}) {
	taskLog := log.(*models.TaskLog)
	session, err := db.GetSession()
	if err == nil {
		defer session.Close()
		db := session.DefaultDB()
		err = db.C("tasklog").Insert(taskLog)
	}

	if err != nil {
		beego.Warn("Insert TaskLog 失败", err)
		//回笼
		t.pushLog(taskLog)
	}
}

func init() {
	if TaskMgt.queue == nil {
		TaskMgt.queue = queue.NewQueue(TaskMgt.saveLog, 100)
	}
}
