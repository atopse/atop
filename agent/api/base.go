package api

import (
	// "github.com/astaxie/beego"

	// "github.com/ysqi/atop/server/core"
	"github.com/ysqi/beegopkg/web"
)

type atopAgentAPIController struct {
	web.APIController
}

func (b *atopAgentAPIController) NeedCheckLogin() bool {
	return false
}
