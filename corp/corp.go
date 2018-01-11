package corp

import (
	"errors"
	"sync"
)

// Corporation 企业微信
type Corporation struct {
	Name           string            `json:"name"`
	ID             string            `json:"corpid"`
	ContactsSecret string            `json:"contacts_secret"`
	AgentMap       map[string]*Agent `json:"agents"`
}

// Agent 应用
type Agent struct {
	Corp       *Corporation
	tokenMutex sync.Mutex
	ID         int    `json:"agentid"`
	Secret     string `json:"secret"`
	MsgToken   string `json:"msgtoken"`
	MsgAesKey  string `json:"msgaeskey"`
}

// GetAgent return *Agent
func (c *Corporation) GetAgent(name string) (*Agent, error) {
	if a, isok := c.AgentMap[name]; isok {
		return a, nil
	}
	return nil, errors.New("agent name not found")
}
