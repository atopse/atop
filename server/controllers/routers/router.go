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

package routers

import (
	"github.com/astaxie/beego"
	"github.com/ysqi/beegopkg/web"

	"github.com/ysqi/atop/server/controllers/api"
	webc "github.com/ysqi/atop/server/controllers/web"
)

func init() {
	web.EnableErrorControll()
	ns := beego.NewNamespace("/api/",
		beego.NSInclude(
			&api.OptCheckController{},
			&api.SystemController{},
			&api.AgentController{},
			&api.MsgController{},
		),
	)
	beego.AddNamespace(ns)
	beego.Include(&webc.CheckConroller{})

	// beego.InsertFilter("/static/*", beego.BeforeStatic, func(ctx *context.Context) {
	// 	ctx.Output.Header("Cache-control", "max-age=5")
	// })
	beego.SetStaticPath("/main.go", "e:/dev/golang/github.com/ysqi/atop/server/main.go")
}
