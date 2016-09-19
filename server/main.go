package main

import (
	"github.com/astaxie/beego"

	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/config"
	"github.com/ysqi/atop/common/log2"

	_ "github.com/ysqi/atop/common/db"
	_ "github.com/ysqi/atop/server/controllers/routers"
)

func main() {
	log2.Info("应用配置文件:", config.AppCfgPath)
	beego.LoadAppConfig("ini", config.AppCfgPath)
	beego.SetLogFuncCall(false)
	beego.AddAPPStartHook(func() error {
		common.RunStartHook()
		return nil
	})
	beego.Run()
}
