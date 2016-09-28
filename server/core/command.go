package core

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ysqi/atop/common/db"
	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/com"
)

var (
	//SuportCmdCategory 所支持的命令类型
	SuportCmdCategory = []string{
		"cmd",
		"ps",
		"mssql",
	}
)

// CommandMgt 命令信息管理实例
var CommandMgt = &CommandService{}

// CommandService 命令信息管理.
type CommandService struct {
}

// NewCommand 保存新的命令
func (c *CommandService) NewCommand(cmd *models.CmdInfo) error {
	cmd.Name = strings.TrimSpace(cmd.Name)
	cmd.Category = strings.TrimSpace(cmd.Category)

	if err := cmd.Verify(); err != nil {
		return err
	}
	if !com.IsSliceContainsStr(SuportCmdCategory, cmd.Category) {
		return errors.New("不支持的命令类型:" + cmd.Category)
	}

	return db.Do(func(d *mgo.Database) error {
		if count, err := d.C("command").Find(bson.M{"name": cmd.Name}).Count(); err != nil {
			return err
		} else if count > 0 {
			return errors.New("已存在相同名称的命令,请修改命令名称")
		}
		cmd.ID = bson.NewObjectId()
		return d.C("command").Insert(cmd)
	})
}

// GetCommandByName 按名称查询Command信息.
func (c *CommandService) GetCommandByName(name string) (cmd *models.CmdInfo, err error) {
	if name == "" {
		return nil, errors.New("命令名称不能为空")
	}
	err = db.FindOne("command", bson.M{"name": name}, cmd)
	return
}
