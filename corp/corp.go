package corp

import (
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
	ID         string `json:"agentid"`
	Secret     string `json:"secret"`
	MsgToken   string `json:"msgtoken"`
	MsgAesKey  string `json:"msgaeskey"`
}

// GetAgent return *Agent
func (c *Corporation) GetAgent(agentName string) *Agent {
	return c.AgentMap[agentName]
}
