package message

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/ysqi/atop/common/log2"
	"github.com/ysqi/atop/common/util"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMessage(t *testing.T) {

	nsqd := []string{"127.0.0.1:4150"}
	nsqlookup := []string{"http://127.0.0.1:4161/", "http://127.0.0.1:4161/"}

	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("atop_test go-nsq/%s", nsq.VERSION)

	Convey("接收消息", t, func() {
		p, err := NewProducer("127.0.0.1:4150", cfg)
		So(err, ShouldBeNil)
		wg := sync.WaitGroup{}

		messages := []interface{}{
			"string", 1, 1.2, "<nil>", false, 'c',
			cfg,
		}

		chanel := "test." + util.NewID()
		d := NewTopicDiscoverer(cfg, chanel, nsqd, nsqlookup)
		h := func(message *nsq.Message) error {
			log2.Debugf("Got Message:%q", string(message.Body))
			wg.Done()
			return nil
		}
		d.SubTopic("test", h)
		go d.Start()
		//自定义Chanel后等待初始化完成
		time.Sleep(2 * time.Second)
		for _, m := range messages {
			wg.Add(1)
			p.Publish("test", m)
		}
		wg.Wait()
		d.Exit()
		p.Exit()
		time.Sleep(5 * time.Second)

	})

}
