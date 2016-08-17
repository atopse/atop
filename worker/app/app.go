package app

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/nsqio/go-nsq"
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

	cfg              config.Configer
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
	if len(topicHandles) == 0 {
		return errors.New("无 Topic Handle 加载,退出程序")
	}

	channel := cfg.DefaultString("channel", "main")

	nsqdTCPAddrs = cfg.Strings("nsqd-tcp-addrs")
	lookupdHTTPAddrs = cfg.Strings("lookupd-http-address")
	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		return errors.New("确实配置项 nsqd-tcp-addrs 或 lookupd-http-address")
	}

	nsqLogger = &NSQLoger{logs.GetBeeLogger()}
	nsqCfg := nsq.NewConfig()
	nsqCfg.UserAgent = fmt.Sprintf("atop_WORKER/%s go-nsq/%s", VERSION, nsq.VERSION)
	nsqCfg.MaxInFlight = 10
	if dialTimeout := time.Duration(cfg.DefaultInt("dialTimeout", 0)); dialTimeout > 0 {
		nsqCfg.DialTimeout = dialTimeout
	}

	discoverer := newTopicDiscoverer(nsqCfg)

	signal.Notify(discoverer.hupChan, syscall.SIGHUP)
	signal.Notify(discoverer.termChan, syscall.SIGINT, syscall.SIGTERM)

	discoverer.watch(channel)
	return nil
}

func init() {
	var err error
	// fmt.Println(os.Args[0])
	cfg, err = config.NewConfig("ini", filepath.Join(os.Args[0], "config", "app.conf"))
	if err != nil {
		panic("解析配置文件 config/app.conf 时失败," + err.Error())
	}
}
