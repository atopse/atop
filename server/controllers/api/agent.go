package api

import (
	"github.com/pquerna/ffjson/ffjson"

	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/server/biz"
	"github.com/ysqi/com"
)

// AgentController Agent相关
type AgentController struct {
	atopAPIController
}

// SayHello  Agent主动跟Server招呼
// @router /agent/sayhello [post]
func (a *AgentController) SayHello() {
	//需要提供Agent信息
	agent := &models.AgentInfo{}
	err := ffjson.Unmarshal(a.Ctx.Input.RequestBody, agent)
	if err != nil {
		a.OutputError(err)
		return
	}
	if err = agent.Verify(); err != nil {
		a.OutputError(err)
		return
	}
	err = biz.AgentMgt.UpdateAgent(*agent, true)
	a.OutputDoResult(err)

}

// Offline Agent 下线通知
// @Param ip body string true "待下线Agent的IP"
// @router /agent/offline [post]
func (a *AgentController) Offline() {
	data := make(map[string]string)
	if err := ffjson.Unmarshal(a.Ctx.Input.RequestBody, &data); err != nil {
		a.OutputError(err)
		return
	}
	ip := data["ip"]
	if ip == "" {
		a.OutputError("缺少参数ip")
		return
	}
	if com.IsIP(ip) == false {
		a.OutputError("参数ip非法")
		return
	}
	err := biz.AgentMgt.UpdateAgentStatus(ip, models.AgentStatusOffline)
	a.OutputDoResult(err)
}
