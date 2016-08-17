package biz

import (
	"errors"
	"io/ioutil"
	"time"

	"git.coding.net/ysqi/atop/common/models"
	"git.coding.net/ysqi/atop/server/src/db"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/toolbox"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ysqi/beegopkg/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
func (o *AgentService) GetAgentInfo(ip string, formDB bool) (*models.AgentInfo, error) {
	if formDB == false {
		agent := o.AgentStore[ip]
		//存在则返回，否则DB从查找
		if agent != nil {
			return agent, nil
		}
	}
	agent := &models.AgentInfo{}
	err := db.Do(func(d *mgo.Database) error {
		return d.C("agent").Find(bson.M{"ip": ip}).One(agent)
	})
	return agent, err
}

// FindAgent 更新Agent状态
func (o *AgentService) FindAgent(agent models.AgentInfo) {
	agent.Status = models.AgentStatusOnline
	if old := o.AgentStore[agent.IP]; old == nil {
		o.AgentStore[agent.IP] = &agent
		err := o.saveAgent(&agent)
		if err != nil {
			beego.Error("保存Agent信息失败,", err)
			return
		}
	} else if old.Status != models.AgentStatusOnline {
		//需要更新数据
		err := o.saveAgent(&agent)
		if err != nil {
			beego.Error("更新Agent信息失败,", err)
			return
		}
		o.AgentStore[agent.IP] = &agent
	}
	if err := o.UpdateAgentStatus(agent.IP, agent.Status); err != nil {
		beego.Warn("更新Agent状态失败，", err)
	}
}

// UpdateAgentStatus 更新 Agent状态
func (o *AgentService) UpdateAgentStatus(agentIP string, status models.AgentStatus) error {
	if agentIP == "" {
		return nil
	}
	agent, err := o.GetAgentInfo(agentIP, false)
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

// GetOnlineAgent 获取Agent
func (o *AgentService) GetOnlineAgent(ip, ip2 string) *models.AgentInfo {
	if ip == "" {
		ip = ip2
	}
	if ip == "" {
		return nil
	}
	agent := o.AgentStore[ip]
	if agent == nil || agent.Status != models.AgentStatusOnline {
		if ip2 == "" {
			return nil
		}
		return o.GetOnlineAgent(ip2, "")
	}
	return agent
}

func (a *AgentService) saveAgent(agent *models.AgentInfo) error {
	if agent == nil {
		return nil
	}
	if agent.IP == "" {
		return errors.New("Agent.IP不能为空")
	}
	return db.Do(func(d *mgo.Database) error {
		c := d.C("agent")
		query := c.Find(bson.M{"ip": agent.IP})
		if count, err := query.Count(); err != nil {
			return err
		} else if count == 0 {
			agent.ID = bson.NewObjectId()
			return c.Insert(agent)
		} else {
			_, err := d.C("agent").Upsert(&bson.M{"ip": agent.IP}, bson.M{
				"$set": bson.M{
					"name":        agent.Name,
					"url":         agent.URL,
					"description": agent.Desc,
					"updated":     time.Now(),
				},
			})
			return err
		}
	})

}
func (o *AgentService) agentHeartbeatChecking() error {
	//从 DB 拉取 Agent
	agents := []*models.AgentInfo{}

	err := db.Do(func(dataBase *mgo.Database) error {
		return dataBase.C("agent").Find(nil).All(&agents)
	})
	if err != nil {
		return err
	}

	//保存到内存
	for _, v := range agents {
		if o.AgentStore[v.IP] == nil {
			o.AgentStore[v.IP] = v
		}
	}

	//循环检查Agent状态
	for _, agent := range o.AgentStore {
		// 未知 或者 在线 时检测
		if agent.Status == models.AgentStatusOnline || agent.Status == models.AgentStatusUnknown {
			u := agent.JionPath("/api/sys/status")
			req := httplib.Get(u)
			result, err := o.doRequest(req)
			if err != nil {
				beego.Warn("检查 Agent 状态失败,", err)
			}
			if result != nil && result.Data == "ok" {
				agent.Status = models.AgentStatusOnline
			} else {
				agent.Status = models.AgentStatusOffline
			}
			if err := o.UpdateAgentStatus(agent.IP, agent.Status); err != nil {
				beego.Warn("更新Agent状态失败，", err)
			}
		}
	}
	return nil
}

// Post POST 提交数据到Agent
func (o *AgentService) Post(agent *models.AgentInfo, path string, data interface{}) (*web.Response, error) {
	u := agent.JionPath("api", path)
	req, err := httplib.Post(u).JSONBody(data)
	if err != nil {
		return nil, err
	}
	return o.doRequest(req)
}

func (o *AgentService) doRequest(req *httplib.BeegoHTTPRequest) (*web.Response, error) {
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
	if err := resData.CheckStatus(); err != nil {
		return resData, err
	}
	return resData, nil
}

func init() {
	//每5分钟执行一次=5*60
	t := toolbox.NewTask("维护更新 Agent", "0/300 * * * * *", func() error {
		err := AgentMgt.agentHeartbeatChecking()
		if err != nil {
			beego.Warn("维护更新 Agent 出现错误:", err)
		}
		return nil
	})
	toolbox.AddTask(t.Taskname, t)
}
