package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.coding.net/ysqi/atop/common/models"
	"git.coding.net/ysqi/atop/server/src/biz"
	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ysqi/beegopkg/web"
	"gopkg.in/mgo.v2/bson"
)

func TestPushCommandMsg(t *testing.T) {
	Convey("测试接收Command进度消息", t, func() {
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
		err := biz.TaskMgt.NewTask(task)
		t.Log(task.ID)
		So(err, ShouldBeNil)
		msg := models.Msg{
			ID:          bson.NewObjectId(),
			Target:      task.ID,
			ContentType: "command",
			Created:     time.Now(),
		}
		process := models.CmdExecProcess{
			CommandID: task.ID,
			Status:    "begin",
			Content: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": 3.33,
				"key4": time.Now(),
			},
		}
		msg.Content = process
		r, _ := http.NewRequest("POST", "/api/msg/command", nil)
		bodyWithJSON(r, msg)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		So(w, ShouldBeGoodResponse)
		actual, err := bufferToStruct(w.Body)
		So(err, ShouldBeNil)
		So(w, ShouldBeGoodResponse)
		So(actual, ShouldBeEqualResponse, &web.Response{Code: 200, Success: true})
	})
}
