package biz

import (
	"errors"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/atop/server/src/db"
	"github.com/astaxie/beego"
)

// OptCheck 三班操作检查项
var OptCheck = &OptCheckService{}

// OptCheckService 三班操作检查项
type OptCheckService struct {
}

var checkItemList []models.CheckItem

func init() {
	beego.AddAPPStartHook(initCheckItemGroup)
}

func initCheckItemGroup() error {
	checkItemList = make([]models.CheckItem, 0, 0)

	session := db.NewSession()
	defer session.Close()

	c := session.DefaultDB().C("checkitem")
	err := c.Find(nil).All(&checkItemList)
	if err != nil {
		return err
	}
	return nil
}

// GetAllCheckItems 获取全部检查项
func (o *OptCheckService) GetAllCheckItems() []models.CheckItem {
	return checkItemList
}

// GetCheckItemByID 根据ID获取检查项
func (o *OptCheckService) GetCheckItemByID(id string) (item *models.CheckItem, err error) {
	session := db.NewSession()
	defer session.Close()
	c := session.DefaultDB().C("checkitem")
	err = c.FindId(id).One(item)
	return
}

// RunOrStopCheckItem 启动或停止检查项
func (o *OptCheckService) RunOrStopCheckItem(id string, goRun bool, operator string) error {
	operator = strings.Trim(operator, "")
	if operator == "" {
		return errors.New("操作人不能为空")
	}
	item, err := o.GetCheckItemByID(id)
	if err != nil {
		return err
	}
	if item.Status == "disabled" {
		return errors.New("该检查项已被禁用，不允许运行")
	}
	if goRun {
		if item.Status == "running" {
			return errors.New("该检查项已在运行中，不能双启动")
		}
	} else {
		if item.Status != "running" {
			return nil
		}
	}
	//运行任务，将此任务下发
	if goRun {
		task, err := o.runCheckItem(item)
		m := bson.M{}
		if task != nil {
			m["options.lastTaskId"] = task.ID
		}
		if err != nil {
			m["status"] = "stoped"
		}
		err2 := db.Do(func(dataBase *mgo.Database) error {
			return dataBase.C("checkitem").UpdateId(item.ID, bson.M{
				"$set": m,
			})
		})
		if err2 != nil {
			beego.Error("更新CheckItem失败,", err2)
		}
		return err
	}
	return nil
}

func (o *OptCheckService) runCheckItem(item *models.CheckItem) (*models.Task, error) {
	task := &models.Task{
		Name:      item.Name,
		TargetIP:  item.TargetIP,
		TargetIP2: item.TargetIP2,
		Options: map[string]interface{}{
			"checkItem": item.ID,
		},
		ResultCheck: item.CheckWays,
		Cmd:         item.Cmd,
	}
	if err := TaskMgt.NewTask(task); err != nil {
		return nil, err
	}
	if err := TaskMgt.StartTask(task); err != nil {
		return task, err
	}
	return task, nil
}
