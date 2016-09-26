package api

// OptCheckController 三班日常操作检查控制
type OptCheckController struct {
	atopAPIController
}

// GetCheckItems 获取检查项列表
// @router /optcheck/list/ [get]
func (o *OptCheckController) GetCheckItems() {
	// items := core.GetAllCheckItems()
	// o.OutputSuccess(items)
}

// RunCheck 执行检查
// @router /optcheck/:checkItemID:string/run [post]
func (o *OptCheckController) RunCheck() {
	// id := o.Ctx.Input.Param(":checkItemID")
}
