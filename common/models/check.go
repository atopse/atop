// Copyright 2016 Author ysqi. All Rights Reserved.
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
// @Author: ysqi
// @Email: devysq@gmail.com or 460857340@qq.com

package models

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/ysqi/atop/common/assertions"
)

// CheckItem 检查项
type CheckItem struct {
	ID        bson.ObjectId     `bson:"_id,omitempty"`
	Name      string            `bson:"name"`       //名称
	TargetIP  string            `bson:"target_ip"`  //执行任务的IP，即在该指定服务器上执行命令
	TargetIP2 string            `bson:"target_ip2"` //代理执行命令的服务器IP
	Cmd       *CmdInfo          `bson:"cmd"`        //命令信息
	Options   map[string]string `bson:"options"`    //更新相关信息
	CheckWays []*ResultCheckWay `bson:"checkways"`  //数据检查方式
	Status    string            `bson:"status"`     //检查项状态：running,stoped,disabled
}

// ResultCheckWay 数据检查模式
type ResultCheckWay struct {
	Way     assertions.Assertion `bson:"way"`
	Params  []interface{}        `bson:"params"`
	Leval   string               `bson:"level"` //Level 基本：info, warn，error
	Options map[string]string    `bson:"options"`
}
