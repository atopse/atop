package main

import (
	"github.com/astaxie/beego"

	_ "git.coding.net/ysqi/atop/server/src/db"
	_ "git.coding.net/ysqi/atop/server/src/routers"
)

func main() {
	beego.Run()
}
