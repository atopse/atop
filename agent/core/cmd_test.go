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

package core

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ysqi/atop/common/models"

	. "github.com/smartystreets/goconvey/convey"
)

var testPort = ":8081"

func TestExecCmd(t *testing.T) {

	err := http.ListenAndServe(testPort, nil)
	if err != nil {
		t.Fatal(err)
	}

	Convey("执行命令测试", t, func() {
		cmd := &models.CmdInfo{
			No:       "001",
			Category: "cmd",
			ResType:  "string",
			Command:  "ipconfig",
		}

		closed := make(chan bool, 1)
		http.HandleFunc(urlSendMsg, func(w http.ResponseWriter, r *http.Request) {
			str, err := ioutil.ReadAll(r.Body)
			ShouldBeNil(err)
			t.Log(str)
		})
		AddNewCmdToPool(cmd)
		<-closed
	})
}
