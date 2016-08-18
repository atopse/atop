package app

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/nsqio/go-nsq"
	"github.com/ysqi/atop/common"
)

// NSQLoger output nsq log info.
type NSQLoger struct {
	log *logs.BeeLogger
}

// Output nsq log with BeeLogger.
func (nlog *NSQLoger) Output(calldepth int, s string) error {
	nlog.log.Write([]byte("[NSQ]" + s))
	return nil
}

func newTopicDiscoverer(cfg *nsq.Config) *TopicDiscoverer {
	return &TopicDiscoverer{
		workers:  make(map[string]*ConsumerWorker),
		termChan: make(chan os.Signal),
		hupChan:  make(chan os.Signal),
		cfg:      cfg,
	}
}

func newConsumerWorker(topic string, channel string, cfg *nsq.Config) (*ConsumerWorker, error) {
	f := &ConsumerWorker{
		Topic:    topic,
		msgChan:  make(chan *nsq.Message, 1),
		ExitChan: make(chan int),
		termChan: make(chan bool),
	}
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}

	consumer.AddHandler(f)
	consumer.SetLogger(nsqLogger, nsq.LogLevelInfo)

	err = consumer.ConnectToNSQDs(nsqdTCPAddrs)
	if err != nil {
		logs.Warn("连接 NSQDS 失败,", err)
	}
	err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs)
	if err != nil {
		logs.Warn("连接 NSQLookupd 失败,", err)
	}

	f.C = consumer
	return f, nil
}

var (
	topicHandles = make(map[string]nsq.HandlerFunc)

	discoverer       *TopicDiscoverer
	nsqLogger        *NSQLoger
	nsqdTCPAddrs     []string
	lookupdHTTPAddrs []string
)

// RegHandle Reg a new topic handle.
// 此说明有可实现的 Handle 来处理 Topic 信息,因此将订阅该主题处理消息
func RegHandle(topic string, h nsq.HandlerFunc) {
	if _, ok := topicHandles[topic]; ok {
		panic(topic + "  已注册,不允许重复")
	}
	if h == nil {
		panic("此 Handle 必需存在")
	}
	topicHandles[topic] = h
}

// Run APP
func Run() error {
	common.RunStartHook()

	if len(topicHandles) == 0 {
		return errors.New("无 Topic Handle 加载,退出程序")
	}

	channel := common.AppCfg.DefaultString("channel", "main")

	nsqdTCPAddrs = common.AppCfg.Strings("nsqd-tcp-addrs")
	lookupdHTTPAddrs = common.AppCfg.Strings("lookupd-http-address")
	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		return errors.New("必需配置 nsqd-tcp-addrs 或 lookupd-http-address")
	}

	nsqLogger = &NSQLoger{logs.GetBeeLogger()}
	nsqCfg := nsq.NewConfig()
	nsqCfg.UserAgent = fmt.Sprintf("atop_WORKER/%s go-nsq/%s", VERSION, nsq.VERSION)
	nsqCfg.MaxInFlight = common.AppCfg.DefaultInt("maxInFlight", 10)
	if dialTimeout := time.Duration(common.AppCfg.DefaultInt("dialTimeout", 0)); dialTimeout > 0 {
		nsqCfg.DialTimeout = dialTimeout
	}

	discoverer = newTopicDiscoverer(nsqCfg)

	signal.Notify(discoverer.hupChan, syscall.SIGHUP)
	signal.Notify(discoverer.termChan, syscall.SIGINT, syscall.SIGTERM)

	discoverer.watch(channel)
	return nil
}

// Stop App
func Stop() {
	if discoverer != nil {
		discoverer.termChan <- syscall.SIGSTOP
	}
}
