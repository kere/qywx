package corp

import (
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
	CorpID      string
	AgentSecret string
}

func newTokenCached(corpID, secret string) *tokenCached {
	t := &tokenCached{CorpID: corpID, AgentSecret: secret}
	t.Init(t)
	return t
}

// CheckValue 检查缓存值是否正确，如果正确缓存
func (t *tokenCached) CheckValue(v interface{}) bool {
	return v.(string) != ""
}

// Build func
func (t *tokenCached) Build(args ...interface{}) (interface{}, int, error) {
	// 获取 access_token
	// 请求方式：GET（HTTPS）
	dat, err := client.Get(fmt.Sprintf(tokenURL, t.CorpID, t.AgentSecret))
	if err != nil {
		return "", 0, err
	}

	// v = newToken(dat.String("access_token"), dat.Int("expires_in"))
	return dat.String("access_token"), dat.Int("expires_in"), nil
}

// map[agentID] *tokenCached
var tokenMap = make(map[int]*tokenCached, 0)

// GetToken get cached token
func (a *Agent) GetToken() (string, error) {
	var t *tokenCached
	var isok bool
	if t, isok = tokenMap[a.ID]; !isok {
		t = newTokenCached(a.Corp.ID, a.Secret)
		tokenMap[a.ID] = t
	}

	tmp := t.Get(a.Corp.ID, a.ID)
	if tmp == nil {
		return "", nil
	}

	return tmp.(string), nil
}
