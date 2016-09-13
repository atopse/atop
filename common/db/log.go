package db

import (
	"github.com/astaxie/beego"
	"github.com/ysqi/atop/common/log2"
	"gopkg.in/mgo.v2"
)

func init() {
	if beego.AppConfig.DefaultBool("dbDebug", false) {
		mgo.SetLogger(log2.GetLogger())
		mgo.SetDebug(true)
	}
}
