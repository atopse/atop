package api

type SystemController struct {
	atopAgentAPIController
}

// Status 获取Agent状态
// @router /sys/status [get]
func (s *SystemController) Status() {
	s.OutputSuccess("OK")
}
