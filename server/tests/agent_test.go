package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/server/core"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ysqi/beegopkg/web"
)

func TestFindAgent(t *testing.T) {

	testCases := []struct {
		title    string
		data     *models.AgentInfo
		expected *web.Response
	}{
		{
			title:    "正常数据测试",
			data:     &models.AgentInfo{Name: "test", URL: "http://127.0.0.1:6060", IP: "127.0.0.1", Desc: "for test"},
			expected: &web.Response{StatusCode: 200},
		},
		{
			title:    "无name测试",
			data:     &models.AgentInfo{Name: "", URL: "http://127.0.0.1:6060", IP: "127.0.0.1", Desc: "for test"},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "无URL测试",
			data:     &models.AgentInfo{Name: "test", URL: "", IP: "127.0.0.1", Desc: "for test"},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "无IP测试",
			data:     &models.AgentInfo{Name: "test", URL: "http://127.0.0.1:6060", IP: "", Desc: "for test"},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "非法URL测试",
			data:     &models.AgentInfo{Name: "test", URL: "http://", IP: "127.0.0.1", Desc: ""},
			expected: &web.Response{StatusCode: 500},
		},
		{
			title:    "非法IP测试",
			data:     &models.AgentInfo{Name: "test", URL: "http://127.0.0.1:6060", IP: "127.0.0", Desc: ""},
			expected: &web.Response{StatusCode: 500},
		},
	}

	Convey("Agent 心跳请求测试", t, func() {
		for _, c := range testCases {
			Convey(c.title, func() {
				actual, err := agentSayHello(c.data)
				So(err, ShouldBeNil)
				So(actual, ShouldBeEqualResponse, c.expected)
			})
		}

	})

}

func TestAgentOffline(t *testing.T) {
	ip := "127.0.0.1"
	type data struct {
		IP string `json:"ip"`
	}
	testCases := []struct {
		title string
		data  data
		want  *web.Response
	}{
		{title: "参数IP不能为空", data: data{IP: ""}, want: &web.Response{StatusCode: 500}},
		{title: "参数IP必须合法", data: data{IP: "127.0.0.555"}, want: &web.Response{StatusCode: 500}},
		{title: "Agent不存在时无法更新", data: data{IP: "127.0.0.111"}, want: &web.Response{StatusCode: 404}},
		{title: "Agent存在时更新成功", data: data{IP: ip}, want: &web.Response{StatusCode: 200}},
	}
	Convey("Agent离线请求测试", t, func() {
		for _, c := range testCases {
			Convey(c.title, func() {
				r, _ := http.NewRequest("POST", "/api/agent/offline", nil)
				bodyWithJSON(r, c.data)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)
				ShouldBeGoodResponse(w)
				actual, err := bufferToStruct(w.Body)
				So(err, ShouldBeNil)
				So(w, ShouldBeGoodResponse)
				So(actual, ShouldBeEqualResponse, c.want)
			})
		}
		Convey("Agent离线DB检查", func() {
			agent, err := core.AgentMgt.GetAgentInfo(ip, true)
			So(err, ShouldBeNil)
			So(agent.IP, ShouldEqual, ip)
			So(agent.Status, ShouldEqual, models.AgentStatusOffline)
		})
	})
}

func agentSayHello(agent *models.AgentInfo) (*web.Response, error) {
	r, _ := http.NewRequest("POST", "/api/agent/sayhello", nil)
	bodyWithJSON(r, agent)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	ShouldBeGoodResponse(w)
	So(w, ShouldBeGoodResponse)
	return bufferToStruct(w.Body)
}
