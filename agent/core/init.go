package core

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ysqi/atop/agent/core/server"
	"github.com/ysqi/atop/common/models"
	"github.com/astaxie/beego"
)

var localhostInfo = &models.AgentInfo{}

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
	if err := server.InitServerInfo(localhostInfo); err != nil {
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
	localhostInfo.Name = cfg["name"]
	localhostInfo.IP = cfg["ip"]

	if localhostInfo.Name == "" {
		localhostInfo.Name, _ = os.Hostname()
	}
	if localhostInfo.IP == "" {
		localhostInfo.IP, _ = externalIP()
	}

	if localhostInfo.Name == "" {
		return errors.New("本机计算机名称尚未配置/获取失败，localhost.name")
	}
	if localhostInfo.IP == "" {
		return errors.New("本机IP尚未配置/获取失败，localhost.ip")
	}
	localhostInfo.URL = fmt.Sprintf("http://%s:%d", localhostInfo.IP, beego.BConfig.Listen.HTTPPort)
	beego.Info(fmt.Sprintf(`
	当前 Agent 信息:
	Name: %s
	IP  : %s
	URL : %s
	`, localhostInfo.Name, localhostInfo.IP, localhostInfo.URL))
	return nil
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// sayHelloWithServer 定期和Server握手
func sayHelloWithServer() {
	period := beego.AppConfig.DefaultInt64("server::heartbeat", 5)
	go func(period int64) {
		var t = time.Duration(period) * time.Minute
		for {

			_, err := server.Post("/agent/sayhello", localhostInfo)

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
