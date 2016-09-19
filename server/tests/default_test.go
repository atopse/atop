package test

import (
	"github.com/ysqi/atop/common"
	"github.com/ysqi/atop/common/config"
	"github.com/ysqi/atop/common/log2"
	_ "github.com/ysqi/atop/server/controllers/routers"

	"github.com/astaxie/beego"
)

func init() {
	log2.Info("应用配置文件:", config.AppCfgPath)
	beego.AddAPPStartHook(func() error {
		common.RunStartHook()
		return nil
	})
	beego.InitBeegoBeforeTest(config.AppCfgPath)
	beego.SetLogFuncCall(false)
}
