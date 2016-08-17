package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ysqi/beegopkg/web"

	"git.coding.net/ysqi/atop/common/models"
)

var serverInfo = models.ServerInfo{}
var currentAgent *models.AgentInfo

func InitServerInfo(agent *models.AgentInfo) error {
	currentAgent = agent

	//加载服务器信息
	c, err := beego.AppConfig.GetSection("server")
	if err != nil {
		return fmt.Errorf("加载服务器信息失败,%s", err)
	}
	if len(c) == 0 {
		return errors.New("请先在配置服务器信息")
	}
	v := c["ip"]
	// if v == "" {
	// 	return errors.New("没有配置服务器IP(server.ip)")
	// }
	serverInfo.IP = v

	v = c["url"]
	if v == "" {
		return errors.New("没有配置服务器URL(server.url)")
	}
	if v[len(v)-1] == '/' {
		v = v[0 : len(v)-1]
	}
	if _, err := url.Parse(v); err != nil {
		return fmt.Errorf("服务器URL(%s)不正确,%s", v, err)
	}
	serverInfo.URL = v

	serverInfo.Name = c["name"]

	beego.Info(fmt.Sprintf(`
	服务器信息：
	server.name = %s
	server.ip	= %s
	server.url	= %s
	`, serverInfo.Name, serverInfo.IP, serverInfo.URL))

	checkServerStatus()

	return nil
}

func checkServerStatus() bool {

	request, err := Get("/sys/ping")
	if err != nil {
		beego.Warn("检查Server 状态失败,", err)
		return false
	}
	if request.Data != "OK" {
		beego.Warn("Server 非正常状态, Status=", request.Data)
		return false
	}
	beego.Info("检查Server", serverInfo.URL, "状态正常")
	return true
}

// // sayHelloWithServer 跟服务器握手.
// // 主动将当前 Agent 信息告知 Server, 让 Server 能感知到有 Agent 活动
// func sayHelloWithServer() error {

// 	u := getServerRequestURL("/agent/sayhello")
// 	req, err := httplib.Post(u.String()).JSONBody(localhostInfo)
// 	if err != nil {
// 		return err
// 	}
// 	res, err := req.DoRequest()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func getServerRequestURL(absURL string) *url.URL {
	u, err := url.Parse(serverInfo.URL)
	if err != nil {
		beego.Error("错误的服务器地址格式,", err)
	}
	u.Path = path.Join(u.Path, "api", absURL)
	return u
}

// Post POST 提交数据到服务器
func Post(path string, data interface{}) (*web.Response, error) {
	u := getServerRequestURL(path)
	req := httplib.Post(u.String()) //.JSONBody(data)

	if data != nil {
		byts, err := ffjson.Marshal(data)
		if err != nil {
			return nil, err
		}
		req.Body(byts)
		req.Header("Content-Type", "application/json")
	}
	return doRequest(req)
}

// Get 提交数据到服务器
func Get(path string, params ...interface{}) (*web.Response, error) {
	u := getServerRequestURL(path)
	req := httplib.Get(u.String())
	return doRequest(req)
}

func doRequest(req *httplib.BeegoHTTPRequest) (*web.Response, error) {
	req.Header("Accept", "application/json")
	response, err := req.DoRequest()
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	resData := &web.Response{}
	if err = ffjson.Unmarshal(body, resData); err != nil {
		return nil, err
	}
	if err = resData.CheckStatus(); err != nil {
		return resData, fmt.Errorf("Server 拒绝请求,%s", err)
	}
	return resData, nil
}
