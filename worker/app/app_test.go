package app

import (
	"testing"

	"github.com/nsqio/go-nsq"
)

func TestAppRun(t *testing.T) {

	RegHandle("test", func(m *nsq.Message) error {
		t.Log(string(m.Body))
		return nil
	})

	if err := Run(); err != nil {
		t.Error(err)
	}

}
