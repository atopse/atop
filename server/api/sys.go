package api

// SystemController 系统相关
type SystemController struct {
	atopAPIController
}

// Ping 获取服务器状态
// @router /sys/ping [get]
func (s *SystemController) Ping() {
	s.OutputSuccess("OK")
}
