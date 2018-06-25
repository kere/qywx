package corp

import (
	"errors"
	"fmt"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/client"
)

const (
	// token 请求URL
	tokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

type tokenCached struct {
	cache.Map
	// CorpID      string
	// AgentSecret string
}

func newTokenCached() *tokenCached {
	t := &tokenCached{}
	t.Init(t)
	return t
}

// CheckValue 检查缓存值是否正确，如果正确缓存
func (t *tokenCached) CheckValue(v interface{}) bool {
	return v.(string) != ""
}

// Build func
func (t *tokenCached) Build(args ...interface{}) (interface{}, int, error) {
	// CorpID: corpID, AgentSecret: secret
	// corpID := args[0].(string)
	corpName := args[0].(string)
	agentID := args[1].(int)
	// cp := GetByID(corpID)
	cp, err := GetByName(corpName)
	if err != nil {
		return nil, 0, errors.New("corp not found in tokenCached")
	}
	agent := cp.GetAgentByAgentid(agentID)
	if agent == nil {
		return nil, 0, errors.New("agent not found in tokenCached")
	}

	// 获取 access_token
	// 请求方式：GET（HTTPS）
	dat, err := client.Get(fmt.Sprintf(tokenURL, cp.Corpid, agent.Secret))
	if err != nil {
		return "", 0, err
	}

	// v = newToken(dat.String("access_token"), dat.Int("expires_in"))
	return dat.String("access_token"), dat.Int("expires_in"), nil
}

// map[agentID] *tokenCached
var tkCached = newTokenCached()

// GetToken get cached token
func (a *Agent) GetToken() (string, error) {
	token := tkCached.Get(a.Corp.Name, a.Agentid)
	if token == nil {
		return "", errors.New("get cached token is nil")
	}

	return token.(string), nil
}
