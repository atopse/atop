package api

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/pquerna/ffjson/ffjson"

	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/server/src/biz"
)

// MsgController 消息处理
type MsgController struct {
	atopAPIController
}

// ReceiveMsg 接收消息，目前支持的消息类型有：1:命令处理进度-command
// @router /msg/:msgType:string [post]
func (m *MsgController) ReceiveMsg() {
	msgType := m.Ctx.Input.Param(":msgType")
	var err error
	switch msgType {
	default:
		err = fmt.Errorf("不支持的消息类型(%s)", msgType)
	case "command":
		err = m.doCommandMsg()
	}
	beego.Debug("error info:", err)
	m.OutputDoResult(err)
}

// doCommandMsg 接收处理命令处理进度消息
func (m *MsgController) doCommandMsg() error {
	msg := &models.Msg{
		Content: &models.CmdExecProcess{},
	}
	if err := ffjson.Unmarshal(m.Ctx.Input.RequestBody, msg); err != nil {
		return err
	}
	if err := msg.Verify(); err != nil {
		return err
	}
	return biz.TaskMgt.PushTaskProcess(msg)
}
