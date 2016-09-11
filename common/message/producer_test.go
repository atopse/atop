package message

import (
	"fmt"
	"testing"

	"sync"

	"github.com/nsqio/go-nsq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPublish(t *testing.T) {
	nsqd := "127.0.0.1:4150"
	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("atop_test/%s go-nsq/%s", "TestPublish", nsq.VERSION)
	Convey("发送消息测试", t, func() {
		p, err := NewProducer(nsqd, cfg)
		So(err, ShouldBeNil)

		Convey("数据非法测试", func() {
			So(func() { p.Publish("", "") }, ShouldPanic)
			So(func() { p.Publish("", "BODY") }, ShouldPanic)
			So(func() { p.Publish("test", nil) }, ShouldPanic)
			So(func() { p.Publish("test", "") }, ShouldPanic)
			So(func() { p.Publish("test", "BODY", "arg") }, ShouldPanic)
		})

		wg := sync.WaitGroup{}
		messages := []interface{}{
			"string", 1, 1.2, "<nil>", false, 'c',
			cfg,
		}
		wg.Add(len(messages))

		after := func(topic string, body interface{}, err error, args ...interface{}) {
			Convey("推送结果检查", t, func() {
				So(topic, ShouldEqual, "test")
				So(body, ShouldBeIn, messages)
				So(err, ShouldBeNil)
				So(args, ShouldContain, "arg1")
				So(args, ShouldContain, "arg2")
				So(args, ShouldContain, "arg3")
			})
			wg.Done()
		}

		for _, m := range messages {
			p.Publish("test", m, after, "arg1", "arg2", "arg3")
		}
		wg.Wait()
		p.Exit()

	})

}
