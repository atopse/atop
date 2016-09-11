package message

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"strconv"

	"github.com/nsqio/go-nsq"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ysqi/atop/common/log2"
)

// AfterMessagePublish 定义在消息完成发送后的自定义事件
type AfterMessagePublish func(topic string, body interface{}, err error, args ...interface{})

type info struct {
	topic string
	body  interface{}
	after AfterMessagePublish
	args  []interface{}
	err   error
}

// ProduceWorker 推送消息队列工作器
type ProduceWorker struct {
	p            *nsq.Producer
	msgChan      chan *info
	afterPublish chan *info
	exitChan     chan bool
	termChan     chan os.Signal
	init         bool
}

func (w *ProduceWorker) router() {
	exist := false
	for {
		select {
		case <-w.exitChan:
			exist = true
		case <-w.termChan:
			exist = true
		case message := <-w.msgChan:
			w.publish(message)
		case message := <-w.afterPublish:
			if message.after == nil {
				if message.err != nil {
					log2.Error(message.err)
				}
			} else {
				//TODO 使用go来执行此事件，不知是否恰当，但至少能防止after执行时间过长影响消息发送.
				go message.after(message.topic, message.body, message.err, message.args...)
			}
		}
		if exist {
			//TODO: 也行还有消息在队列中待发送
			w.p.Stop()
			close(w.termChan)
			close(w.msgChan)
			close(w.exitChan)
			break
		}
	}

}

func (w *ProduceWorker) publish(message *info) {
	if message == nil {
		return
	}
	var byteBody []byte
	switch v := message.body.(type) {
	case []byte:
		byteBody = v
	case string:
		byteBody = []byte(v)
	case bool:
		byteBody = strconv.AppendBool(byteBody, v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		byteBody = []byte(fmt.Sprintf("%#v", v))
	default:
		byteBody, message.err = ffjson.Marshal(message.body)
		if message.err != nil {
			message.err = fmt.Errorf("推送消息编译为JSON时出错,body=%#v,err=%s", message.body, message.err)
			w.afterPublish <- message
			return
		}
		// 移除前后的引号""
		if len(byteBody) > 2 && byteBody[0] == '"' && byteBody[len(byteBody)-1] == '"' {
			byteBody = byteBody[1 : len(byteBody)-1]
		}
	}

	message.err = w.p.Publish(message.topic, byteBody)
	if message.err != nil {
		message.err = fmt.Errorf("推送消息失败,body=%q,err=%s", string(byteBody), message.err)
	}
	w.afterPublish <- message
}

// Exit 退出ProduceWorker.
func (w *ProduceWorker) Exit() {
	w.exitChan <- true
}

// Publish 推送消息到Topic.
// topic 和 message 参数不能为空，如果为空，则忽略.
// 同时 Message 对象将转换为JSON格式发送，除非 Message 类型为[]byte.
func (w *ProduceWorker) Publish(topic string, message interface{}, afterWithargs ...interface{}) {
	if w.init == false {
		panic("使用错误，尚未正常初始化ProduceWorker")
	}
	if topic == "" {
		panic("topic不能为空")
	}
	if message == nil {
		panic("message不能为空")
	}
	if str, ok := message.(string); ok && str == "" {
		panic("message不能为空")
	}
	item := &info{
		topic: topic,
		body:  message,
	}
	if len(afterWithargs) > 0 {
		after, ok := afterWithargs[0].(func(topic string, body interface{}, err error, args ...interface{}))
		if !ok {
			panic("第一个参数必须是 func(topic string, body interface{}, err error, args ...interface{} ")
		}
		item.after = after
		item.args = append(item.args, afterWithargs[1:]...)
	}
	w.msgChan <- item
}

// NewProducer 新建推送队列工作器.
// 创建的实例将自动维护状态,如果cfg为空，则创建默认config.
func NewProducer(addr string, cfg *nsq.Config) (*ProduceWorker, error) {
	if cfg == nil {
		cfg = nsq.NewConfig()
		cfg.UserAgent = fmt.Sprintf("atop go-nsq/%s", nsq.VERSION)
	}
	p, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		return nil, err
	}
	worker := &ProduceWorker{
		p:            p,
		msgChan:      make(chan *info, 10000),
		afterPublish: make(chan *info, 10000),
		exitChan:     make(chan bool),
		termChan:     make(chan os.Signal),
	}
	signal.Notify(worker.termChan, syscall.SIGINT, syscall.SIGTERM)
	go worker.router()
	worker.init = true
	return worker, nil
}
