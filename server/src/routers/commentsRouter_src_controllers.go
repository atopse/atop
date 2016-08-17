package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["git.coding.net/ysqi/atop/server/src/controllers:CheckConroller"] = append(beego.GlobalControllerRouter["git.coding.net/ysqi/atop/server/src/controllers:CheckConroller"],
		beego.ControllerComments{
			"Index",
			`/check/:t:string`,
			[]string{"get"},
			nil})

}
