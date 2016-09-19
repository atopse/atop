package main

import (
	"github.com/astaxie/beego"

	_ "github.com/ysqi/atop/server/db"
	_ "github.com/ysqi/atop/server/routers"
)

func main() {
	beego.Run()
}
