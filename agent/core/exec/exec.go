package exec

import (
	"errors"
	"fmt"

	"github.com/ysqi/atop/common/models"
)

var commds = map[string]Commd{}

// Commd 借款信息
type Commd interface {
	Exec(cmd *models.CmdInfo) (result interface{}, err error)
}

// Reg 注册
func Reg(apater string, cmd Commd) error {
	if apater == "" {
		return errors.New("注册名不能为空")
	}
	if _, ok := commds[apater]; ok {
		return fmt.Errorf("%s已注册，不允许重复", apater)
	}
	commds[apater] = cmd
	return nil
}

// Run 执行命令
func Run(cmd *models.CmdInfo) (result interface{}, err error) {
	c, ok := commds[cmd.Category]
	if !ok {
		return nil, fmt.Errorf("不支持的命令类型：%s", cmd.Category)
	}
	return c.Exec(cmd)
}
