package app

import (
	nsq "github.com/nsqio/go-nsq"
)

// ConsumerWorker 消费者 Worker
type ConsumerWorker struct {
	Topic    string
	C        *nsq.Consumer
	msgChan  chan *nsq.Message
	ExitChan chan int
	termChan chan bool
	hupChan  chan bool
}

// HandleMessage 处理消息
func (w *ConsumerWorker) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	w.msgChan <- m
	return nil
}

func (w *ConsumerWorker) router(r *nsq.Consumer) {
	closing := false
	exit := false

	msgHandle := topicHandles[w.Topic].HandleMessage
	for {
		select {
		case <-r.StopChan:
			exit = true
		case <-w.termChan:
			closing = true
		case m := <-w.msgChan:
			//消息处理
			if err := msgHandle(m); err != nil {
				m.Requeue(0)
			} else {
				m.Finish()
			}
		}
		if closing {
			r.Stop()
		}
		if exit {
			r.Stop()
			close(w.ExitChan)
			break
		}
	}
}
