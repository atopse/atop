package db

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
)

type log struct {
	beegoLog *logs.BeeLogger
}

func (l *log) Output(calldepth int, s string) error {
	l.beegoLog.Debug("%s", s)
	return nil
}

func init() {
	if beego.AppConfig.DefaultBool("dbDebug", false) {
		logExtend := &log{
			beegoLog: beego.BeeLogger,
		}
		mgo.SetLogger(logExtend)
		mgo.SetDebug(true)
	}
}
