package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ysqi/atop/common/models"
	_ "github.com/ysqi/atop/server/routers"

	"github.com/astaxie/beego"
	"github.com/pquerna/ffjson/ffjson"
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
			expected: &web.Response{Code: 200, Success: true},
		},
		{
			title:    "无name测试",
			data:     &models.AgentInfo{Name: "", URL: "http://127.0.0.1:6060", IP: "127.0.0.1", Desc: "for test"},
			expected: &web.Response{Code: 5000, Success: false},
		},
		{
			title:    "无URL测试",
			data:     &models.AgentInfo{Name: "test", URL: "", IP: "127.0.0.1", Desc: "for test"},
			expected: &web.Response{Code: 5000, Success: false},
		},
		{
			title:    "无IP测试",
			data:     &models.AgentInfo{Name: "test", URL: "http://127.0.0.1:6060", IP: "", Desc: "for test"},
			expected: &web.Response{Code: 5000, Success: false},
		},
		{
			title:    "非法URL测试",
			data:     &models.AgentInfo{Name: "test", URL: "http://", IP: "127.0.0.1", Desc: ""},
			expected: &web.Response{Code: 5000, Success: false},
		},
		{
			title:    "非法IP测试",
			data:     &models.AgentInfo{Name: "test", URL: "http://127.0.0.1:6060", IP: "127.0.0", Desc: ""},
			expected: &web.Response{Code: 5000, Success: false},
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
func ShouldBeGoodResponse(actual interface{}, expected ...interface{}) string {
	res := actual.(*httptest.ResponseRecorder)
	if res == nil {
		return "结果为<nil>"
	}
	if res.Code != 200 {
		return fmt.Sprintf("HTTP Code 期望是200，实际上是：%d", res.Code)
	}
	exp := "application/json; charset=utf-8"
	if t := res.HeaderMap.Get("Content-Type"); t != exp {
		return fmt.Sprintf("Response Content-Type 期望是%s,实际上是%s", exp, t)
	}
	return ""
}

func ShouldBeEqualResponse(actual interface{}, expected ...interface{}) string {
	result := actual.(*web.Response)
	exp := expected[0].(*web.Response)
	if exp == nil {
		return "expected=<nil>"
	}
	if exp.Code != result.Code && result.Code != 200 {
		printJSON(result)
	}
	if exp.Code != result.Code {
		return fmt.Sprintf("Response Code 期望是%d，实际上是：%d", exp.Code, result.Code)
	}
	if exp.Success != result.Success {
		return fmt.Sprintf("Response Success 期望是%v，实际上是：%v", exp.Success, result.Success)
	}
	return ""
}

func printJSON(data interface{}) {
	content, err := ffjson.Marshal(data)
	if err != nil {
		return
	}
	fmt.Println(string(content))
}
func bodyWithJSON(req *http.Request, data interface{}) error {

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if data == nil {
		return nil
	}
	byts, err := ffjson.Marshal(data)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(byts))
	req.ContentLength = int64(len(byts))
	return nil
}
func bufferToStruct(buffer *bytes.Buffer) (*web.Response, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return nil, err
	}
	resData := &web.Response{}
	if err = ffjson.Unmarshal(body, resData); err != nil {
		beego.Debug("Respone Body:\n", string(body), "\n ------body end")
		return nil, err
	}
	return resData, nil
}
