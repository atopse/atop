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

package core

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
)

var configPath = getConfigPath()

func getConfigPath() string {
	_, file, _, _ := runtime.Caller(1)
	dir := filepath.Dir(file)
	for {
		d := filepath.Join(dir, "conf", "app.test.conf")
		if dir == filepath.VolumeName(dir) {
			return d
		}
		if utils.FileExists(d) {
			return d
		}
		// Parent dir.
		dir = filepath.Dir(dir)
	}
}

func TestGetIP(t *testing.T) {
	ip, err := externalIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ip)
}

func init() {
	beego.BConfig.RunMode = "test"
	if err := beego.LoadAppConfig("ini", configPath); err != nil {
		panic(err)
	}

	if err := initInfo(); err != nil {
		panic(err)
	}
}
