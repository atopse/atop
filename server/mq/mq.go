package mq

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

func initMQ() error {
	msCfg, err := beego.AppConfig.GetSection("mq")
	if err != nil {
		return err
	}
	url := msCfg["url"]
	if url == "" {
		return errors.New("尚未配置MQ.url 信息")
	}
	config := nsq.NewConfig()
	producer, err = nsq.NewProducer(url, config)
	if err != nil {
		return err
	}
	return nil
}
