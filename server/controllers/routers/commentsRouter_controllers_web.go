package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/web:CheckConroller"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers/web:CheckConroller"],
		beego.ControllerComments{
			Method: "Index",
			Router: `/check/:t:string`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
