package app

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"

	"log"

	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/config"
	"github.com/ysqi/atop/common/message"
)

var (
	topicHandles     = make(map[string]nsq.HandlerFunc)
	discoverer       *message.TopicDiscoverer
	nsqdTCPAddr      string
	lookupdHTTPAddrs []string
	nsqCfg           *nsq.Config
)
var (
	// ProduceWorker 用于发送信息到 NSQD
	ProduceWorker *message.ProduceWorker
)

// AddHandle Add a new topic handle.
// 此说明有可实现的 Handle 来处理 Topic 信息,因此将订阅该主题处理消息
func AddHandle(topic string, h nsq.HandlerFunc) error {
	return discoverer.SubTopic(topic, h)
}
func init() {

	channel := config.AppCfg.DefaultString("channel", "main")

	nsqdTCPAddr = config.AppCfg.String("nsqd-tcp-addr")
	lookupdHTTPAddrs = config.AppCfg.Strings("lookupd-http-address")
	if nsqdTCPAddr == "" {
		log.Fatalln("必需配置 nsqd-tcp-addrs")
	}
	if len(lookupdHTTPAddrs) == 0 {
		log.Fatalln("必需配置 lookupd-http-address")
	}

	nsqCfg = nsq.NewConfig()
	nsqCfg.UserAgent = fmt.Sprintf("atop_WORKER/%s go-nsq/%s", VERSION, nsq.VERSION)
	nsqCfg.MaxInFlight = config.AppCfg.DefaultInt("maxInFlight", 10)
	if dialTimeout := time.Duration(config.AppCfg.DefaultInt("dialTimeout", 0)); dialTimeout > 0 {
		nsqCfg.DialTimeout = dialTimeout
	}
	discoverer = message.NewTopicDiscoverer(nsqCfg, channel, nil, lookupdHTTPAddrs)
}

// Run APP
func Run() error {
	common.RunStartHook()
	var err error
	ProduceWorker, err = message.NewProducer(nsqdTCPAddr, nsqCfg)
	if err != nil {
		return err
	}
	go discoverer.Start()

	return nil
}

// Stop App
func Stop() {
	if discoverer != nil {
		discoverer.Exit()
	}
	if ProduceWorker != nil {
		ProduceWorker.Exit()
	}
}
