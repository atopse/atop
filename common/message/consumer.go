package message

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	nsq "github.com/nsqio/go-nsq"
	"github.com/ysqi/atop/common/log2"
)

// ConsumerWorker 消费者 Worker
type ConsumerWorker struct {
	topic            string
	C                *nsq.Consumer
	ExitChan         chan bool
	msgChan          chan *nsq.Message
	handler          func(topic string, message *nsq.Message)
	termChan         chan bool
	channel          string   //topic所在channel
	nsqdTCPAddrs     []string //NSQD地址
	lookupdHTTPAddrs []string //Lookupd地址
}

// HandleMessage 处理消息,用于从NSD接收消息
func (w *ConsumerWorker) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	w.msgChan <- m
	return nil
}

func (w *ConsumerWorker) watch() {
	exit := false
	ticker := time.Tick(5 * time.Minute)
	w.connectSync()

	for {
		select {
		case <-ticker:
			w.connectSync()
		case <-w.C.StopChan:
			exit = true
		case <-w.termChan:
			exit = true
		case m := <-w.msgChan:
			w.handler(w.topic, m)
		}
		if exit {
			w.C.Stop()
			//TODO: 如果此时msgChan中还有缓存队列消除存在时，则会导致消息丢失，是否应该编译msgChan清理消息?

			close(w.msgChan)
			close(w.termChan)
			w.ExitChan <- true
			break
		}
	}
}
func (w *ConsumerWorker) connectSync() {
	err := w.C.ConnectToNSQLookupds(w.lookupdHTTPAddrs)
	if err != nil {
		logs.Warn("连接 NSQLookupd 失败,", err)
	}
	err = w.C.ConnectToNSQDs(w.nsqdTCPAddrs)
	if err != nil {
		log2.Warn("连接 NSQDS 失败,", err)
	}
	log2.Debugf("连接NSQDS/NSQLookupd，status=%#v", w.C.Stats())
}

// NewConsumerWorker 创建消费者对象
func NewConsumerWorker(nsqdTCPAddrs []string, lookupdHTTPAddrs []string, topic string, channel string, handler func(topic string, message *nsq.Message), cfg *nsq.Config) (*ConsumerWorker, error) {
	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		return nil, errors.New("NSQDS和NSQLookupd不能全为空")
	}
	if handler == nil {
		return nil, errors.New("消息处理函数不能为空")
	}
	if topic == "" {
		return nil, errors.New("topic不能为空")
	}
	if channel == "" {
		return nil, errors.New("channel不能为空")
	}
	w := &ConsumerWorker{
		topic:            topic,
		msgChan:          make(chan *nsq.Message, 1),
		ExitChan:         make(chan bool),
		termChan:         make(chan bool),
		handler:          handler,
		channel:          channel,
		nsqdTCPAddrs:     nsqdTCPAddrs,
		lookupdHTTPAddrs: lookupdHTTPAddrs,
	}
	consumer, err := nsq.NewConsumer(w.topic, w.channel, cfg)
	if err != nil {
		return nil, err
	}
	w.C = consumer
	consumer.AddHandler(w)

	// TODO: 更加全局日志基本设置
	consumer.SetLogger(log2.GetLogger(), nsq.LogLevelDebug)

	return w, nil
}
