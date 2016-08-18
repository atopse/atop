package main

import (
	"os"

	"github.com/ysqi/atop/worker/app"

	"github.com/astaxie/beego/logs"
)

func main() {
	if err := app.Run(); err != nil {
		logs.Error("初始化失败,", err)
		os.Exit(1)
	}
}
