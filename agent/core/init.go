package core

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/ysqi/atop/agent/core/server"
	"github.com/ysqi/atop/common/models"
	"github.com/ysqi/com"
)

// CurrentAgent 当前Ageng实例信息
var CurrentAgent = &models.AgentInfo{}

func init() {

	if beego.BConfig.RunMode != "test" {
		beego.AddAPPStartHook(initInfo)
	}
}

func initInfo() error {
	if err := initLocalhostInfo(); err != nil {
		return err
	}
	beego.Info("初始化 Agent 信息完成")
	if err := server.InitServerInfo(CurrentAgent); err != nil {
		return err
	}
	beego.Info("初始化 Server 信息完成")
	sayHelloWithServer()
	beego.Info("初始化心跳任务完成")
	return nil
}

// initLocalhostInfo 初始化当前 Agent 信息.
// 优先从配置文件读取配置信息,如果尚未配置则获取当前运行 Agent 服务器的计算机名称和可用 IP 作为默认配置.
func initLocalhostInfo() error {
	cfg, err := beego.AppConfig.GetSection("localhost")
	if err != nil {
		return err
	}
	CurrentAgent.Name = cfg["name"]
	CurrentAgent.IP = cfg["ip"]

	if CurrentAgent.Name == "" {
		CurrentAgent.Name, _ = os.Hostname()
	}
	if CurrentAgent.IP == "" {
		CurrentAgent.IP, err = com.ExternalIP()
		if err != nil {
			panic(err)
		}
	}

	if CurrentAgent.Name == "" {
		return errors.New("本机计算机名称尚未配置/获取失败，localhost.name")
	}
	if CurrentAgent.IP == "" {
		return errors.New("本机IP尚未配置/获取失败，localhost.ip")
	}
	CurrentAgent.URL = fmt.Sprintf("http://%s:%d", CurrentAgent.IP, beego.BConfig.Listen.HTTPPort)
	beego.Info(fmt.Sprintf(`
	当前 Agent 信息:
	Name: %s
	IP  : %s
	URL : %s
	`, CurrentAgent.Name, CurrentAgent.IP, CurrentAgent.URL))
	return nil
}

// sayHelloWithServer 定期和Server握手
func sayHelloWithServer() {
	period := beego.AppConfig.DefaultInt64("server::heartbeat", 5)
	go func(period int64) {
		var t = time.Duration(period) * time.Minute
		for {

			_, err := server.Post("/agent/sayhello", CurrentAgent)

			//如果请求超时,则10秒后重试
			if err == http.ErrHandlerTimeout {
				time.Sleep(10 * time.Second)
				continue
			}
			if err != nil {
				beego.Warn("心跳发送失败,", err)
			} else {
				beego.Debug("心跳请求成功")
			}

			if err != nil && beego.BConfig.RunMode == "dev" {
				time.Sleep(5 * time.Second)
				continue
			}
			time.Sleep(t)
		}
	}(period)

	//TODO:需支持退出 Agent 时,主动告知 Server 已下线
}
