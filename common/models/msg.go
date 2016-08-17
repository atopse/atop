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
	"time"

	"github.com/ysqi/beegopkg/models"

	"gopkg.in/mgo.v2/bson"
)

// Msg 消息
type Msg struct {
	ID          bson.ObjectId //消息唯一ID
	Target      bson.ObjectId //消息主题ID
	ContentType string        //消息内容类型，方便进行数据解析
	Content     interface{}   //消息内容
	SendTimes   int           //发送次数，表示该消息被发送过的次数
	Created     time.Time     //消息创建时间
}

// Verify 验证数据合法性
func (m *Msg) Verify() error {
	if m.ID == "" {
		return &models.VerifyError{"消息ID不能为空"}
	}
	if m.Target == "" {
		return &models.VerifyError{"消息TargetID不能为空"}
	}
	if m.ContentType == "" {
		return &models.VerifyError{"消息ContentType不能为空"}
	}
	if m.Content == nil {
		return &models.VerifyError{"消息Content不能为空"}
	}
	return nil
}
