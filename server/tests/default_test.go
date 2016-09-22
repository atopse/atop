package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/config"
	"github.com/ysqi/atop/common/log2"
	"github.com/ysqi/beegopkg/web"

	"strconv"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ysqi/atop/server/controllers/api"
	_ "github.com/ysqi/atop/server/controllers/routers"
)

func init() {
	log2.Info("应用配置文件:", config.AppCfgPath)
	beego.AddAPPStartHook(func() error {
		common.RunStartHook()
		return nil
	})
	beego.InitBeegoBeforeTest(config.AppCfgPath)
	beego.SetLogFuncCall(false)
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
	if exp.StatusCode != result.StatusCode && result.StatusCode != 200 {
		printJSON(result)
	}
	if exp.StatusCode != result.StatusCode {
		return fmt.Sprintf("Response Code 期望是%d，实际上是：%d", exp.StatusCode, result.StatusCode)
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

func TestUnmarshalBody(t *testing.T) {
	Convey("解析请求包", t, func() {

		type dt struct {
			Name  string `json:"Name"`
			Value int64  `json:"value"`
		}
		var data = dt{
			Name:  "ysqi",
			Value: 12,
		}

		check := func(r *http.Request) {
			c := api.CmdController{}
			c.Ctx = context.NewContext()
			c.Ctx.Reset(nil, r)
			c.Ctx.Input.CopyBody(1024)
			result := &dt{}
			err := c.UnmarshalBody(result)
			So(err, ShouldBeNil)
			So(result.Name, ShouldEqual, data.Name)
			So(result.Value, ShouldEqual, data.Value)
		}
		Convey("JSON请求", func() {
			r, _ := http.NewRequest("POST", "/api", nil)
			bodyWithJSON(r, data)
			check(r)
		})
		Convey("Form表单请求POST", func() {
			r, _ := http.NewRequest("POST", "/api", nil)
			r.Form = url.Values{}
			r.Form.Add("Name", data.Name)
			r.Form.Add("Value", strconv.FormatInt(data.Value, 10))
			check(r)
		})
		Convey("Form表单请求GET", func() {
			r, _ := http.NewRequest("GET", "/api", nil)
			r.Form = url.Values{}
			r.Form.Add("Name", data.Name)
			r.Form.Add("Value", strconv.FormatInt(data.Value, 10))
			check(r)
		})
	})
}
