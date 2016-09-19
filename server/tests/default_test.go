package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/config"
	"github.com/ysqi/atop/common/log2"
	_ "github.com/ysqi/atop/server/controllers/routers"
	"github.com/ysqi/beegopkg/web"

	"github.com/astaxie/beego"
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
