package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:AgentController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:AgentController"],
		beego.ControllerComments{
			Method: "SayHello",
			Router: `/agent/sayhello`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:AgentController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:AgentController"],
		beego.ControllerComments{
			Method: "Offline",
			Router: `/agent/offline`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:MsgController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:MsgController"],
		beego.ControllerComments{
			Method: "ReceiveMsg",
			Router: `/msg/:msgType:string`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:OptCheckController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:OptCheckController"],
		beego.ControllerComments{
			Method: "GetCheckItems",
			Router: `/optcheck/list/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:OptCheckController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:OptCheckController"],
		beego.ControllerComments{
			Method: "RunCheck",
			Router: `/optcheck/:checkItemID:string/run`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:SystemController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/api:SystemController"],
		beego.ControllerComments{
			Method: "Ping",
			Router: `/sys/ping`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
