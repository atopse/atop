package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ysqi/atop/common/assertions"
	"github.com/ysqi/atop/common/log2"
	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/beegopkg/web"
	"github.com/ysqi/com"
)

func unmarshalData(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func TestWhereJSON(t *testing.T) {
	byteDate := []byte(`{"taskID":"abc"}`)

	var data interface{}

	if err := json.Unmarshal(byteDate, &data); err != nil {
		t.Fatal(err)
	}
	t.Logf("Success ,taskID=%s", data.(map[string]interface{})["taskID"])

	f := func(data []byte, v interface{}) error {
		err := json.Unmarshal(data, v)
		return err
	}
	if err := f(byteDate, &data); err != nil {
		t.Fatal(err)
	}
	t.Logf("Success ,taskID=%s", data.(map[string]interface{})["taskID"])

	if err := unmarshalData(byteDate, &data); err != nil {
		t.Fatal(err)
	}
	t.Logf("Success ,taskID=%s", data.(map[string]interface{})["taskID"])

}
func TestNewTask(t *testing.T) {
	cmd := &models.CmdInfo{
		Name:     "检查网络配置",
		Category: "cmd",
		ResType:  "string",
		Command:  "ipconfig",
	}
	cwItems := []*models.ResultCheckWay{&models.ResultCheckWay{Way: assertions.ShouldBeEmpty}}

	testCases := []struct {
		title    string
		data     *models.Task
		expected *web.Response
	}{
		{
			title:    "Name不能未空",
			data:     &models.Task{TargetIP: "127.0.0.1", Cmd: cmd, ResultCheck: cwItems},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "TargetIP不能未空",
			data:     &models.Task{Name: "Test Task", Cmd: cmd, ResultCheck: cwItems},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "Command不能未空",
			data:     &models.Task{Name: "Test Task", TargetIP: "127.0.0.1", ResultCheck: cwItems},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "ResultCheck不能未空",
			data:     &models.Task{Name: "Test Task", TargetIP: "127.0.0.1", Cmd: cmd},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "TargetIP不允许非法",
			data:     &models.Task{Name: "Test Task", TargetIP: "127.0.0.1.abc", Cmd: cmd, ResultCheck: cwItems},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "TargetIP2不允许非法",
			data:     &models.Task{Name: "Test Task", TargetIP: "127.0.0.1", TargetIP2: "127.0.0.1.abc", Cmd: cmd, ResultCheck: cwItems},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "正常数据",
			data:     &models.Task{Name: "Test Task", TargetIP: "127.0.0.1", TargetIP2: "127.0.0.2", Cmd: cmd, ResultCheck: cwItems},
			expected: &web.Response{StatusCode: 200},
		},
	}
	Convey("新建任务", t, func() {
		for _, c := range testCases {
			Convey(c.title, func() {
				r, _ := http.NewRequest("POST", "/api/task/new", nil)
				bodyWithJSON(r, c.data)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)
				actual, err := bufferToStruct(w.Body)
				So(err, ShouldBeNil)
				So(w, ShouldBeGoodResponse)
				So(actual, ShouldBeEqualResponse, c.expected)
			})
		}

	})
}

func TestStartTask(t *testing.T) {

	// agent := &models.AgentInfo{Name: "test", URL: "http://127.0.0.1:6060", IP: "127.0.0.1", Desc: "for test"}

	task := newTask()
	testCases := []struct {
		title    string
		data     interface{}
		expected *web.Response
	}{
		{title: "任务ID不允许为空", data: map[string]interface{}{"taskID": ""}, expected: &web.Response{StatusCode: 500}},
		{title: "任务ID不允许非法", data: "12345", expected: &web.Response{StatusCode: 500}},
		{title: "任务不存在", data: map[string]interface{}{"taskID": bson.NewObjectId().Hex()}, expected: &web.Response{StatusCode: 404}},
		{title: "任务正常启动", data: map[string]interface{}{"taskID": task.ID.Hex()}, expected: &web.Response{StatusCode: 200}},
		{title: "任务不允许重复执行", data: map[string]interface{}{"taskID": task.ID.Hex()}, expected: &web.Response{StatusCode: 500}},
	}

	time.Sleep(10000)
	Convey("启动任务", t, func() {
		for _, c := range testCases {
			Convey(c.title, func() {
				r, _ := http.NewRequest("POST", "/api/task/start", nil)
				bodyWithJSON(r, c.data)
				w := httptest.NewRecorder() 
				beego.BeeApp.Handlers.ServeHTTP(w, r)
				actual, err := bufferToStruct(w.Body)
				So(err, ShouldBeNil)
				So(w, ShouldBeGoodResponse)
				So(actual, ShouldBeEqualResponse, c.expected)
			})
		}

	})
}

func newTask() *models.Task {
	cmd := &models.CmdInfo{
		Name:     "检查网络配置",
		Category: "cmd",
		ResType:  "string",
	}
	if runtime.GOOS == "windows" {
		cmd.Command = "echo %GOROOT%"
	} else {
		cmd.Command = "echo $GOROOT"
	}
	cwItems := []*models.ResultCheckWay{
		&models.ResultCheckWay{
			Way:   assertions.ShouldEqual,
			Leval: "error",
			Params: map[string]interface{}{
				"want": os.Getenv("GOROOT"),
			},
		},
	}
	ip, err := com.ExternalIP()
	if err != nil {
		log2.Fatalln(err)
	}
	task := &models.Task{
		Name:        "Test Task",
		TargetIP:    ip,
		TargetIP2:   "127.0.0.1",
		Cmd:         cmd,
		ResultCheck: cwItems}
	r, _ := http.NewRequest("POST", "/api/task/new", nil)
	bodyWithJSON(r, task)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	actual, err := bufferToStruct(w.Body)
	if err != nil {
		log2.Fatalln(err)
	}
	task.ID = bson.ObjectIdHex(actual.Data.(map[string]interface{})["taskID"].(string))
	return task
}
