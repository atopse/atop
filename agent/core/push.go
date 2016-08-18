package core

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	"github.com/otium/queue"

	"github.com/ysqi/atop/agent/core/server"
	"github.com/ysqi/atop/common/models"
)

var msgCenter = &MsgCenter{}

// MsgCenter 消息推送处理
type MsgCenter struct {
	queue *queue.Queue
	once  sync.Once
}

func (m *MsgCenter) initQueue() {
	m.queue = queue.NewQueue(func(msg interface{}) {
		m.sendToServer(msg.(*models.Msg))
	}, 100)
}

// PushMsg 推送消息到服务器
func PushMsg(msgType string, msg interface{}) {
	msgCenter.once.Do(msgCenter.initQueue)
	msgCenter.push(&models.Msg{
		ID:          bson.NewObjectId(),
		Content:     msg,
		ContentType: msgType,
		Created:     time.Now(),
	})
}

// push 消息入队列
func (m *MsgCenter) push(msg *models.Msg) {
	if msg.Created.IsZero() {
		msg.Created = time.Now()
	}
	m.queue.Push(msg)
}

// 发送消息失败重复次数
const reTrySendMaxTimes = 3

// sendToServer 发送消息到服务器
func (m *MsgCenter) sendToServer(msg *models.Msg) {
	_, err := server.Post("/msg/"+msg.ContentType, msg)
	if err == nil {
		return
	}
	if msg.SendTimes > reTrySendMaxTimes {
		beego.Warn(fmt.Sprintf("发送数据%#v到服务器失败:%s", msg.Content, err))
		return
	}
	beego.Warn(fmt.Sprintf("[重试]发送数据%#v到服务器失败:%s", msg.Content, err))
	// TODO: 每次重发，应该等待一段时间

	//重新放入队列发送
	msg.SendTimes++
	m.push(msg)
}
