// Copyright 2016 The Ysqi Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");

//自动化运维 Server 端程序 .
//
//基于 Beego 框架实现的Go Web 程序,作为自动化运维软件的Server端,负责和各Agent通信以及提供操作端UI.
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
