package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:AgentController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:AgentController"],
		beego.ControllerComments{
			"SayHello",
			`/agent/sayhello`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:AgentController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:AgentController"],
		beego.ControllerComments{
			"Offline",
			`/agent/offline`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:MsgController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:MsgController"],
		beego.ControllerComments{
			"ReceiveMsg",
			`/msg/:msgType:string`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:OptCheckController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:OptCheckController"],
		beego.ControllerComments{
			"GetCheckItems",
			`/optcheck/list/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:OptCheckController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:OptCheckController"],
		beego.ControllerComments{
			"RunCheck",
			`/optcheck/:checkItemID:string/run`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:SystemController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/src/api:SystemController"],
		beego.ControllerComments{
			"Status",
			`/sys/status`,
			[]string{"get"},
			nil})

}
