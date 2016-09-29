package core

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/toolbox"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ysqi/beegopkg/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ysqi/atop/common/config"
	"github.com/ysqi/atop/common/db"
	"github.com/ysqi/atop/common/log2"
	"github.com/ysqi/atop/common/models"
)

// AgentMgt Agent管理
var AgentMgt = &AgentService{
	AgentStore: make(map[string]*models.AgentInfo),
}

// AgentService Agent管理器
type AgentService struct {
	AgentStore map[string]*models.AgentInfo
}

// GetAgentInfo 获取Agent信息，可选择从数据库获取，否则从内存中获取.
//
// 如果 fromDB 为 true,则从数据库中根据IP查询匹配的Agent信息,否则直接从内存中获取.
// 但是当内存中不存在此Ageng信息时将强制从数据库中查找.
func (as *AgentService) GetAgentInfo(ip string, formDB bool) (*models.AgentInfo, error) {
	if formDB == false {
		agent := as.AgentStore[ip]
		//存在则返回，否则DB从查找
		if agent != nil {
			return agent, nil
		}
	}
	agent := &models.AgentInfo{}
	err := db.Do(func(d *mgo.Database) error {
		return d.C("agent").Find(bson.M{"ip": ip}).One(agent)
	})
	//如果从数据库中有查找到Agent则更新到内存中.
	if err != nil && agent.ID != "" {
		as.AgentStore[ip] = agent
	}
	return agent, err
}

// UpdateAgent 依据IP更新Agent,同步更新内存和数据库.
// 更新数据库时仅仅更新:name,url,description,status,updated字段.
func (as *AgentService) UpdateAgent(agent models.AgentInfo, online bool) error {
	if err := agent.Verify(); err != nil {
		return err
	}
	if online {
		agent.Status = models.AgentStatusOnline
	} else {
		if a := as.AgentStore[agent.IP]; a != nil {
			agent.Status = a.Status
		} else {
			agent.Status = models.AgentStatusUnknown
		}
	}
	//优先保存到内存
	as.AgentStore[agent.IP] = &agent

	return db.Do(func(d *mgo.Database) error {
		c := d.C("agent")
		if count, err := c.Find(bson.M{"ip": agent.IP}).Count(); err != nil {
			return err
		} else if count == 0 {
			agent.ID = bson.NewObjectId()
			return c.Insert(agent)
		} else {
			_, err := c.Upsert(&bson.M{"ip": agent.IP}, bson.M{
				"$set": bson.M{
					"name":        agent.Name,
					"url":         agent.URL,
					"description": agent.Desc,
					"status":      agent.Status,
					"updated":     time.Now(),
				},
			})
			return err
		}
	})
}

// UpdateAgentStatus 更新 Agent状态.
func (as *AgentService) UpdateAgentStatus(ip string, status models.AgentStatus) error {
	if ip == "" {
		return nil
	}
	agent, err := as.GetAgentInfo(ip, false)
	if err != nil {
		return err
	}
	agent.Status = status
	agent.Updated = time.Now()

	// 更新到数据库
	return db.Do(func(d *mgo.Database) error {
		return d.C("agent").Update(bson.M{"ip": agent.IP}, bson.M{
			"$set": bson.M{
				"status":  agent.Status,
				"updated": agent.Updated,
			},
		})
	})
}

// GetOnlineAgent 返回在线的Agent.
// 优先根据ip检查内存中在线的Agent,否则依次从backupIP 列表中检查,直到得到在线的Agent.
func (as *AgentService) GetOnlineAgent(ip string, backupIPs ...string) *models.AgentInfo {
	agent := as.AgentStore[ip]
	if agent == nil || agent.Status != models.AgentStatusOnline {
		// 发送请求到Agent,检查是否在线
		response, err := as.HTTPDoRequest(
			&models.AgentInfo{IP: ip, URL: fmt.Sprintf("http://%s:%d", ip, config.AppCfg.DefaultInt("agentPort", 8909))},
			"get", "ping",
		)
		if err == nil && response.Data != nil {
			if a, ok := response.Data.(map[string]interface{}); ok {
				agent = &models.AgentInfo{
					IP:   a["IP"].(string),
					URL:  a["URL"].(string),
					Name: a["Name"].(string),
				}
				log2.Debugf("ping agent<%s> success,the default url is %s", agent.String(), agent.URL)

				if err = as.UpdateAgent(*agent, true); err != nil {
					log2.Warnf("主动更新Agent信息失败,", err)
				} else {
					return agent
				}
			}

		}

		if len(backupIPs) == 0 {
			return nil
		}
		return as.GetOnlineAgent(backupIPs[0], backupIPs[1:]...)
	}
	return agent
}

// HTTPDoRequest 用HTTP提交数据到指定Agent.
// 返回请求结果数据
func (as *AgentService) HTTPDoRequest(agent *models.AgentInfo, method string, path string, data ...interface{}) (*web.Response, error) {
	u := agent.JionPath("api", path)
	var req *httplib.BeegoHTTPRequest

	if method == "post" {
		req = httplib.Post(u)
		if len(data) > 0 {
			var err error
			req, err = req.JSONBody(data[0])
			if err != nil {
				return nil, err
			}
		}
	} else {
		req = httplib.Get(u)
	}
	return as.doRequest(req)
}

func (as *AgentService) doRequest(req *httplib.BeegoHTTPRequest) (*web.Response, error) {
	req.Header("Accept", "application/json")
	req.Header("Agent", fmt.Sprintf("ATOPServer/%s", config.AppCfg.String("version")))
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
	if err := resData.CheckStatus(); err != nil {
		return resData, err
	}
	return resData, nil
}

// agentHeartbeatChecking 依次将数据库所保存的Agent发送心跳包.
//
// 心跳发送前会将内存数据同数据库同步一次,主要是补充内存中尚未保存的Agent信息.
// 注意心跳包只会发送给非离线的Agent.
func (as *AgentService) agentHeartbeatChecking() {
	//从 DB 拉取 Agent
	agents := []*models.AgentInfo{}

	err := db.Do(func(dataBase *mgo.Database) error {
		return dataBase.C("agent").Find(nil).All(&agents)
	})
	if err != nil {
		log2.Warn("心跳前获取数据库Agent列表失败,", err.Error())
	} else {
		//保存到内存
		for _, v := range agents {
			if as.AgentStore[v.IP] == nil {
				as.AgentStore[v.IP] = v
			}
		}
	}

	//循环检查Agent状态
	for _, agent := range as.AgentStore {
		// 未知 或者 在线 时检测
		if agent.Status == models.AgentStatusOnline || agent.Status == models.AgentStatusUnknown {
			result, err := as.HTTPDoRequest(agent, "get", "sys/ping")
			if err != nil {
				log2.Debugf("检查 Agent<%s> 状态失败,%s", agent.String(), err)
			}
			if result != nil && result.Data == "ok" {
				agent.Status = models.AgentStatusOnline
			} else {
				agent.Status = models.AgentStatusOffline
			}
			if err := as.UpdateAgentStatus(agent.IP, agent.Status); err != nil {
				log2.Debugf("更新 Agent<%s> 状态失败,%s", agent.String(), err)
			}
		}
	}
}

func init() {
	// BUG(ysqi)  将时间周期设置为可配置.
	//每5分钟执行一次=5*60
	t := toolbox.NewTask("维护更新 Agent", "0/300 * * * * *", func() error {
		AgentMgt.agentHeartbeatChecking()
		return nil
	})
	toolbox.AddTask(t.Taskname, t)
}
