package api

import (
	// "github.com/astaxie/beego"

	// "github.com/ysqi/atop/server/src/biz"
	"github.com/ysqi/beegopkg/web"
)

type atopAgentAPIController struct {
	web.APIController
}

func (b *atopAgentAPIController) NeedCheckLogin() bool {
	return false
}
