package exec

import (
	"bytes"
	"os/exec"

	"github.com/ysqi/com"

	"git.coding.net/ysqi/atop/common/models"
	"github.com/henrylee2cn/mahonia"
)

// OSCommand OS command 执行
type OSCommand struct {
}

// Exec 执行命令
func (c *OSCommand) Exec(cmd *models.CmdInfo) (interface{}, error) {
	args := []string{"/C", cmd.Command}
	args = append(args, cmd.Args...)
	command := exec.Command("cmd", args...)
	randomBytes := &bytes.Buffer{}
	command.Stdout = randomBytes
	if err := command.Start(); err != nil {
		return nil, err
	}
	var needWait = true
	if cmd.Options != nil {
		sync, err := com.StrTo(cmd.Options["sync"]).Bool()
		if err != nil {
			return nil, err
		}
		needWait = sync
	}
	if needWait {
		if err := command.Wait(); err != nil {
			return nil, err
		}
	}
	data := string(randomBytes.Bytes())
	// TODO: 性能优化，如果操作系统非中文则不需要进行转码
	return mahonia.NewDecoder("gb2312").ConvertString(data), nil
}

func init() {
	Reg("cmd", &OSCommand{})
}
