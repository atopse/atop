package controllers

import "github.com/ysqi/atop/agent/core"

type SystemController struct {
	baseController
}

// Status 获取Agent状态
// @router /ping [get]
func (s *SystemController) Ping() {
	s.OutputSuccess(core.CurrentAgent)
}
