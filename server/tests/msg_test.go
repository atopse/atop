package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/server/biz"
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
			Content: models.CmdExecProcess{
				CommandID: task.ID,
				Tag:       "processing",
				Body: map[string]interface{}{
					"key1": "value1",
					"key2": 2,
					"key3": 3.33,
					"key4": time.Now(),
				},
			},
		}
		r, _ := http.NewRequest("POST", "/api/msg/command", nil)
		bodyWithJSON(r, msg)
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)
		So(w, ShouldBeGoodResponse)
		actual, err := bufferToStruct(w.Body)
		So(err, ShouldBeNil)
		So(w, ShouldBeGoodResponse)
		So(actual, ShouldBeEqualResponse, &web.Response{StatusCode: 200})
	})
}
