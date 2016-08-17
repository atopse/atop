package api

import (
	// "github.com/astaxie/beego"

	// "git.coding.net/ysqi/atop/server/src/biz"
	"github.com/ysqi/beegopkg/web"
)

type atopAgentAPIController struct {
	web.APIController
}

func (b *atopAgentAPIController) NeedCheckLogin() bool {
	return false
}
