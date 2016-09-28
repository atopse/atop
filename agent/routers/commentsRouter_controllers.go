package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ysqi/atop/agent/controllers:CmdController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/agent/controllers:CmdController"],
		beego.ControllerComments{
			Method: "Exec",
			Router: `/command/exec`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/ysqi/atop/agent/controllers:SystemController"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/agent/controllers:SystemController"],
		beego.ControllerComments{
			Method: "Ping",
			Router: `/ping`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
