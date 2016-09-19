package api

import (
	"github.com/pquerna/ffjson/ffjson"

	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/server/biz"
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
	biz.AgentMgt.FindAgent(*agent)
	a.OutputSuccess()
}

// Offline Agent 下线通知
// @router /agent/offline [post]
func (a *AgentController) Offline() {
	ip := a.GetString("ip")
	if ip == "" {
		a.OutputError("缺少参数ip")
		return
	}
	biz.AgentMgt.UpdateAgentStatus(ip, models.AgentStatusOffline)
	a.OutputSuccess()
}
