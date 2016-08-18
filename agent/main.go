package main

import (
	_ "github.com/ysqi/atop/agent/docs"
	_ "github.com/ysqi/atop/agent/routers"

	"github.com/astaxie/beego"
)

//bee run -gendoc=true -downdoc=true

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
