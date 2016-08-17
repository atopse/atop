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
	"fmt"
	"net/url"
	"path"

	"gopkg.in/mgo.v2/bson"

	"time"

	"github.com/ysqi/com"
)

// ServerInfo 服务器信息
type ServerInfo struct {
	Name    string `bson:"name"`    //名称
	IP      string `bson:"ip"`      //IP
	URL     string `bson:"url"`     //URL
	Version string `bson:"version"` //版本
}

// AgentInfo Agent信息
type AgentInfo struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Name    string        `bson:"name"`        //名称
	IP      string        `bson:"ip"`          //IP
	URL     string        `bson:"url"`         //URL
	Version string        `bson:"version"`     //版本
	Status  AgentStatus   `bson:"status"`      //状态
	Desc    string        `bson:"description"` //描述
	Updated time.Time     `bson:"updated" json:"-"`
}

func (a *AgentInfo) String() string {
	return a.Name + "(" + a.IP + ")"
}

// Verify 验证
func (a *AgentInfo) Verify() error {
	if a.Name == "" {
		return errors.New("名称为空")
	}
	if a.IP == "" {
		return errors.New("IP未空")
	}
	if a.URL == "" {
		return errors.New("URL未空")
	}
	//非法性验证
	if com.IsUrl(a.URL) == false {
		return fmt.Errorf("非法URL：%s", a.URL)
	}
	if com.IsIP(a.IP) == false {
		return fmt.Errorf("非法IP：%s", a.IP)
	}
	return nil
}

// JionPath 拼接 URL
func (a *AgentInfo) JionPath(ps ...string) string {
	u, _ := url.Parse(a.URL)
	elems := []string{u.Path}
	u.Path = path.Join(append(elems, ps...)...)
	return u.String()
}

// AgentStatus Agent状态
type AgentStatus string

const (
	AgentStatusUnknown AgentStatus = ""        //未知
	AgentStatusOnline  AgentStatus = "online"  //在线
	AgentStatusOffline AgentStatus = "offline" //离线

)
