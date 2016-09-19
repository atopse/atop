package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers:CheckConroller"] = append(beego.GlobalControllerRouter["github.com/ysqi/atop/server/controllers:CheckConroller"],
		beego.ControllerComments{
			"Index",
			`/check/:t:string`,
			[]string{"get"},
			nil})

}
