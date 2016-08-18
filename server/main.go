package main

import (
	"github.com/astaxie/beego"

	_ "github.com/ysqi/atop/server/src/db"
	_ "github.com/ysqi/atop/server/src/routers"
)

func main() {
	beego.Run()
}
