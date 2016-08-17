package core

import (
	"errors"
	"fmt"

	"github.com/otium/queue"

	"git.coding.net/ysqi/atop/agent/core/exec"
	m "git.coding.net/ysqi/atop/common/models"
)

// AddNewCmdToPool Add a new cmd to command pool.
func AddNewCmdToPool(cmd *m.CmdInfo) error {
	if err := cmd.Verify(); err != nil {
		return err
	}
	// Push to queue.
	cmdQueue.Push(cmd)
	return nil
}

// cmd exec queue.
var cmdQueue *queue.Queue

const execCMDConcurrencyLimit = 20

func init() {
	cmdQueue = queue.NewQueue(func(cmd interface{}) {
		execCmd(cmd.(*m.CmdInfo))
	}, execCMDConcurrencyLimit)
}

// ExecCmdTask 添加命令任务
func ExecCmdTask(cmdInfo *m.CmdInfo) error {
	if cmdInfo == nil {
		return errors.New("命令不能为空")
	}
	if err := cmdInfo.Verify(); err != nil {
		return err
	}
	cmdQueue.Push(cmdInfo)
	return nil
}

// execCmd 执行命令并反馈处理进度
// 根据命令执行返回信息,信息包括状态信息,执行结果信息,不同信息类型,使用不同说明.
// 执行过程:
// 		1. 状态报告:processing  处理中
//		2. 开始执行命令,如果命令执行失败,则返回: error ,并停止任务
//		3. 执行完成, 返回 result 结果信息
//		4. 状态报告:执行完成 completed
func execCmd(cmdInfo *m.CmdInfo) {
	//status: started,processing,stopped,completed
	info := m.CmdExecProcess{
		CommandID: cmdInfo.ID,
	}
	newProcess := func(newStatus string) {
		info.Content.Tag = "newStatus"
		info.Content.Body = newStatus
		PushMsg("command", info)
	}
	sendErrorStop := func(err interface{}) {
		info.Content.Tag = "error"
		info.Content.Body = fmt.Sprintf("异常，命令已停止,%v", err)
		PushMsg("command", info)
	}
	defer func() {
		err := recover()
		if err != nil {
			sendErrorStop(err)
		}
	}()
	newProcess("processing")
	result, err := exec.Run(cmdInfo)
	if err != nil {
		sendErrorStop(err)
		return
	}
	info.Content.Tag = "result"
	info.Content.Body = result
	PushMsg("command", info)
	newProcess("completed")
}
