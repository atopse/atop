package biz

import (
	"testing"

	"git.coding.net/ysqi/atop/common/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRunNewTask(t *testing.T) {
	Convey("执行Task", t, func() {
		task := &models.Task{
			Name:      "ForTest",
			TargetIP:  "127.0.0.1",
			TargetIP2: "",
			Cmd: &models.CmdInfo{
				Category: "cmd",
				ResType:  "string",
				Command:  "ipconfig",
			},
		}
		err := TaskMgt.NewTask(task)
		So(err, ShouldBeNil)
		err = TaskMgt.StartTask(task)
		So(err, ShouldNotBeNil)
		t.Log(err)
	})
}
