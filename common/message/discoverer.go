package message

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
	nsq "github.com/nsqio/go-nsq"
	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/log2"
)

// TopicDiscoverer Topic 集中管理
type TopicDiscoverer struct {
	workers          map[string]*ConsumerWorker //一个topic对应一个消费者worker
	topicHandles     map[string]nsq.HandlerFunc //每个topic消息对应的消息处理函数
	channel          string                     //topic所在channel
	nsqdTCPAddrs     []string                   //NSQD地址
	lookupdHTTPAddrs []string                   //Lookupd地址
	exitChan         chan bool
	termChan         chan os.Signal
	wg               sync.WaitGroup
	cfg              *nsq.Config
	locker           *sync.Mutex
}

func (td *TopicDiscoverer) startTopicWatch(c *ConsumerWorker) {
	log2.Debugf("开始监控Consumer,status=%v", c.C.Stats())
	td.wg.Add(1)
	defer td.wg.Done()
	go c.watch()
	<-c.ExitChan
}

func (td *TopicDiscoverer) handleMessage(topic string, message *nsq.Message) {
	handle := td.topicHandles[topic]
	if handle == nil {
		log2.Errorf("收到来自主题<%s>的消息<%s>Attempts=%d,但无法处理", topic, message.ID, message.Attempts)
		message.Requeue(0)
		return
	}
	//消息处理
	if err := handle(message); err != nil {
		if badE, ok := err.(common.ErrBadBody); ok {
			logs.Warn("[拒绝处理]消息ID=%s,%s", message.ID, badE.Error())
		} else {
			logs.Warn("[重新提交]消息ID=%s,%s", err.Error())
			message.Requeue(0)
			return
		}
	}
	message.Finish()
}

func (td *TopicDiscoverer) syncTopics() {
	td.locker.Lock()
	for topic := range td.topicHandles {
		if _, exist := td.workers[topic]; exist {
			continue
		}
		consumer, err := NewConsumerWorker(td.nsqdTCPAddrs, td.lookupdHTTPAddrs, topic, td.channel, td.handleMessage, td.cfg)
		if err != nil {
			log2.Warnf("创建 Topic Handle Worker %q 失败,channel=%q,%s", topic, td.channel, err)
			continue
		}
		if _, exist := td.workers[topic]; exist {
			consumer.C.Stop()
			continue
		}
		td.workers[topic] = consumer
		go td.startTopicWatch(consumer)
	}
	td.locker.Unlock()
}

// SubTopic 订阅指定Topic.
// 订阅时不允许重复订阅，一个Topic只能被一个Handle接收并处理.
func (td *TopicDiscoverer) SubTopic(topic string, handle nsq.HandlerFunc) error {
	if topic == "" {
		return errors.New("Topic不能为空")
	}
	if _, ok := td.topicHandles[topic]; ok {
		return fmt.Errorf("该Topic<%s>已被订阅,不运行重复处理", topic)
	}
	td.topicHandles[topic] = handle
	return nil
}

// Start 开启
func (td *TopicDiscoverer) Start() {
	signal.Notify(td.termChan, syscall.SIGINT, syscall.SIGTERM)
	td.syncTopics()
	ticker := time.Tick(5 * time.Minute)
	var exist bool
	for {
		select {
		case <-ticker:
			td.syncTopics()
		case <-td.exitChan:
			exist = true
		case <-td.termChan:
			exist = true
		}
		if exist {
			td.exit()
			break
		}
	}
}

func (td *TopicDiscoverer) exit() {
	for _, w := range td.workers {
		w.termChan <- true
	}
	close(td.exitChan)
	close(td.termChan)

	td.wg.Wait()
}

// Exit 退出
func (td *TopicDiscoverer) Exit() {
	td.exitChan <- true
	td.wg.Wait()
}

// NewTopicDiscoverer 初始化创建主题订阅器.
func NewTopicDiscoverer(cfg *nsq.Config, channel string, nsqdTCPAddrs, lookupdHTTPAddrs []string) *TopicDiscoverer {
	return &TopicDiscoverer{
		workers:          make(map[string]*ConsumerWorker),
		termChan:         make(chan os.Signal),
		exitChan:         make(chan bool),
		topicHandles:     make(map[string]nsq.HandlerFunc),
		cfg:              cfg,
		locker:           &sync.Mutex{},
		channel:          channel,
		nsqdTCPAddrs:     nsqdTCPAddrs,
		lookupdHTTPAddrs: lookupdHTTPAddrs,
	}
}
