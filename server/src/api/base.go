package api

import (
	// "github.com/astaxie/beego"

	// "github.com/ysqi/atop/server/src/biz"
	"github.com/ysqi/beegopkg/web"
)

type atopAPIController struct {
	web.APIController
}

func (b *atopAPIController) NeedCheckLogin() bool {
	return false
}
