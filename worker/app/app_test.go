package app

import (
	"testing"

	"github.com/nsqio/go-nsq"
)

func TestAppRun(t *testing.T) {
	topics := []string{"test", "t2", "t3", "t4"}
	for _, tp := range topics {
		AddHandle(tp, func(m *nsq.Message) error {
			return nil
		})
	}
	if err := Run(); err != nil {
		t.Error(err)
	} else {
		Stop()
	}
}
