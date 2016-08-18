package task

import (
	"github.com/nsqio/go-nsq"
	"github.com/ysqi/atop/worker/app"
)

// HandleNewStatus will update task status.
func HandleNewStatus(m *nsq.Message) error {
    
	return nil
}

func init() {
	app.RegHandle("task.newstatus", HandleNewStatus)
}
