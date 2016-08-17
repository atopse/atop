package app

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	nsq "github.com/nsqio/go-nsq"
)

// TopicDiscoverer Topic 集中管理
type TopicDiscoverer struct {
	workers  map[string]*ConsumerWorker
	termChan chan os.Signal
	hupChan  chan os.Signal
	wg       sync.WaitGroup
	cfg      *nsq.Config
}

func (t *TopicDiscoverer) startTopicRouter(c *ConsumerWorker) {
	t.wg.Add(1)
	defer t.wg.Done()
	go c.router(c.C)
	<-c.ExitChan
}

func (t *TopicDiscoverer) syncTopics(channel string) {
	for topic := range topicHandles {
		if _, exist := t.workers[topic]; exist {
			continue
		}
		consumer, err := newConsumerWorker(topic, channel, t.cfg)
		if err != nil {
			logs.Warn("创建 Topic Handle Worker %q 失败, %s", topic, err)
			continue
		}
		t.workers[topic] = consumer
		go t.startTopicRouter(consumer)
	}
}

func (t *TopicDiscoverer) stop() {
	for _, w := range t.workers {
		w.termChan <- true
	}
}

func (t *TopicDiscoverer) watch(channel string) {
	t.syncTopics(channel)
	ticker := time.Tick(5 * time.Minute)
	for {
		select {
		case <-ticker:
			log.Println("syncTopics...")
			t.syncTopics(channel)
		case <-t.termChan:
			t.stop()
			t.wg.Wait()
			return
		}
	}
}
