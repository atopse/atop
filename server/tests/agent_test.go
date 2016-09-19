package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ysqi/atop/common/models"

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
				r, _ := http.NewRequest("POST", "/api/agent/sayhello", nil)
				bodyWithJSON(r, c.data)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)
				ShouldBeGoodResponse(w)
				actual, err := bufferToStruct(w.Body)
				So(err, ShouldBeNil)
				So(w, ShouldBeGoodResponse)
				So(actual, ShouldBeEqualResponse, c.expected)

				//检查数据库数据

			})
		}

	})

}
