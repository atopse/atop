// Copyright 2016 Author Think. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// @Author: Think
// @Email: devysq@gmail.com or 460857340@qq.com

package models

import (
	"errors"
	"strings"
	"time"

	"github.com/ysqi/com"
	"gopkg.in/mgo.v2/bson"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusNew        TaskStatus = "new"        //新任务
	TaskStatusStarted    TaskStatus = "started"    // 已发送
	TaskStatusProcessing TaskStatus = "processing" // 执行中
	TaskStatusErrorDown  TaskStatus = "errorDown"  // 错误停止
	TaskStatusCompleted  TaskStatus = "completed"  // 已完成
)

// 任务状态集合
var TaskStatusItems = [...]TaskStatus{
	TaskStatusNew,
	TaskStatusStarted,
	TaskStatusProcessing,
	TaskStatusErrorDown,
	TaskStatusCompleted,
}

// Task 任务
type Task struct {
	ID          bson.ObjectId          `bson:"_id,omitempty"`
	Name        string                 `bson:"name"`        //任务名称
	TargetIP    string                 `bson:"targetIp"`    //执行任务的IP，即在该指定服务器上执行命令
	TargetIP2   string                 `bson:"targetIp2"`   //代理执行命令的服务器IP
	Status      TaskStatus             `bson:"status"`      //状态
	Operator    string                 `bson:"operator"`    //操作者
	Options     map[string]interface{} `bson:"options"`     //更新相关信息
	ResultCheck []*ResultCheckWay      `bson:"resultCheck"` //数据检查方式
	Cmd         *CmdInfo               `bson:"cmd"`         //命令信息
}

// TaskLog 任务日志
type TaskLog struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	TaskID         bson.ObjectId `bson:"taskId"`
	Content        interface{}   `bson:"content"`
	OccurrenceTime time.Time     `bson:"occurrenceTime"`
}

// CmdExecProcess CMD命令执行进度报告
type CmdExecProcess struct {
	CommandID bson.ObjectId `bson:"_id,omitempty"` //对应的 command ID
	Tag       string        `bson:"tag"`
	Body      interface{}   `bson:"body"`
	// Status string //状态：started,processing,stopped,completed
	// Process   int           //进度
}

// CmdResultType 命令执行结果数据类型格式
type CmdResultType string

const (
	ResultTypeIsSingleNumber CmdResultType = "number" // 单个数字
	ResultTypeIsSingleString CmdResultType = "string" // 单个字符
	ResultTypeIsOneRow       CmdResultType = "single_row"
	ResultTypeIsMultiRows    CmdResultType = "multi_rows"
	ResultTypeIsJSON         CmdResultType = "json"
)

// CmdInfo 命令信息
type CmdInfo struct {
	ID       bson.ObjectId     `bson:"_id,omitempty"`
	Category string            `bson:"category"` //命令类型： cmd, ps, sql
	ResType  CmdResultType     `bson:"resType"`  //命令结果类型：数字：number, 文本：string, 多行记录：records
	Command  string            `bson:"command"`  //命令内容
	Options  map[string]string `bson:"options"`  //更多命令执行选项信息
	Args     []string          `bson:"args"`
}

var (
	//SuportCmdCategory 所支持的命令类型
	SuportCmdCategory = []string{
		"cmd",
		"ps",
		"mssql",
	}
)

// Verify 验证输入
func (c *CmdInfo) Verify() error {

	if c.Category == "" {
		return errors.New("命令类型不能为空")
	}
	if c.ResType == "" {
		return errors.New("命令执行结果数据类型不能为空")
	}
	if c.Command == "" {
		return errors.New("命令内容不能为空")
	}

	c.Category = strings.ToLower(c.Category)

	if !com.IsSliceContainsStr(SuportCmdCategory, c.Category) {
		return errors.New("不支持的命令类型:" + c.Category)
	}

	return nil
}
