package corp

import (
	"errors"
)

// Corporation 企业微信
type Corporation struct {
	ID             int               `json:"id"`
	Corpid         string            `json:"corpid"`
	Name           string            `json:"name"`
	Title          string            `json:"title"`
	ContactsSecret string            `json:"contacts_secret"`
	ContactsToken  string            `json:"contacts_token"`
	ContactsAesKey string            `json:"contacts_aeskey"`
	MainURL        string            `json:"main_url"`
	AgentMap       map[string]*Agent `json:"agents"`
}

// Agent 应用
type Agent struct {
	Corp      *Corporation
	ID        int    `json:"id"`
	Agentid   int    `json:"agentid"`
	App       int    `json:"app"` //判断应用类型
	ParentID  int    `json:"parent_id"`
	Parent    *Agent `json:"-"` //JSON时忽略字段
	Name      string `json:"name"`
	Secret    string `json:"secret"`
	MsgToken  string `json:"msgtoken"`
	MsgAesKey string `json:"msgaeskey"`
}

// GetAgent return *Agent
func (c *Corporation) GetAgent(name string) (*Agent, error) {
	if a, isok := c.AgentMap[name]; isok {
		return a, nil
	}
	return nil, errors.New("agent name not found")
}

// GetAgentByAgentidID return *Agent
func (c *Corporation) GetAgentByAgentid(id int) *Agent {
	for _, a := range c.AgentMap {
		if a.Agentid == id {
			return a
		}
	}
	return nil
}

// GetAgentByID return *Agent
func (c *Corporation) GetAgentByID(id int) *Agent {
	for _, a := range c.AgentMap {
		if a.ID == id {
			return a
		}
	}
	return nil
}
