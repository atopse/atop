package app

import (
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
)

func TestAppRun(t *testing.T) {
	topics := []string{"test", "t2", "t3", "t4"}
	for _, tp := range topics {
		RegHandle(tp, func(m *nsq.Message) error {
			t.Log(tp + ":" + string(m.Body))
			return nil
		})
	}
	exitC := make(chan bool, 1)
	go func() {
		timer := time.NewTimer(2 * time.Second)
		for {
			select {
			case <-timer.C:
				if len(discoverer.workers) != len(topics) {
					t.Fatalf("lenght of discoverer.workers  want %d,got %d", len(topics), len(discoverer.workers))
				}
				Stop()
				return
			case <-exitC:
				return
			}

		}
	}()
	if err := Run(); err != nil {
		t.Error(err)
		exitC <- true
	}
}
