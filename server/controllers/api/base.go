package api

import (
	// "github.com/astaxie/beego"

	// "github.com/ysqi/atop/server/core"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/ysqi/beegopkg/web"
)

type atopAPIController struct {
	web.APIController
}

func (b *atopAPIController) NeedCheckLogin() bool {
	return false
}

// UnmarshalBody 将请求数据解析为对象.
// 如果请求包为JSON格式数据则使用ffjson包解析对象,否则使用ParseForm解析,并返回错误信息.
func (b *atopAPIController) UnmarshalBody(value interface{}) error {
	if b.Ctx.Input.AcceptsJSON() {
		if err := ffjson.Unmarshal(b.Ctx.Input.RequestBody, value); err != nil {
			return err
		}
	}
	return b.ParseForm(value)
}
