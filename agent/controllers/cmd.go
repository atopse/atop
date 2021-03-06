package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"

	"github.com/ysqi/atop/agent/core"
	m "github.com/ysqi/atop/common/models"
)

// CmdController 接受并处理命令
type CmdController struct {
	baseController
}

// @Title ExecCommand
// @Description  get a command and goto exec.
// @Param	body		body 	models.Response	true   CmdInfo
// @Success 200 {string} models.Response
// @Failure 403 body is empty
// @router /command/exec [post]
func (c *CmdController) Exec() {
	var cmd *m.CmdInfo
	beego.Debug(string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cmd)
	if err != nil {
		c.OutputError(err)
		return
	}
	err = core.ExecCmdTask(cmd)
	c.OutputDoResult(err)
}
