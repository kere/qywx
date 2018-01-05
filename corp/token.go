package corp

import (
	"fmt"
	"time"

	"github.com/kere/qywx/client"
)

const (
	// token 请求URL
	tokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

// Token corp access token
type Token struct {
	Value     string
	Expires   int
	ExpiresAt time.Time
}

// newToken build token
func newToken(token string, expires int) *Token {
	t := &Token{Value: token}
	t.Expires = expires
	t.ExpiresAt = time.Now().Add(time.Duration(expires-60) * time.Second)
	return t
}

var tokenMap = make(map[string]*Token, 0)

// GetToken get cached token
func (a *Agent) GetToken() (*Token, error) {
	now := time.Now()
	key := a.Corp.ID + a.ID

	a.tokenMutex.Lock()

	t, isok := tokenMap[key]
	if isok || now.Before(t.ExpiresAt) {
		a.tokenMutex.Unlock()
		return t, nil
	}

	// 获取 access_token
	// 请求方式：GET（HTTPS）
	uri := fmt.Sprintf(tokenURL, a.Corp.ID, a.Secret)
	dat, err := client.GetMapData(uri)
	if err != nil {
		a.tokenMutex.Unlock()
		return nil, err
	}

	t = newToken(dat.String("access_token"), dat.Int("expires_in"))
	tokenMap[key] = t
	a.tokenMutex.Unlock()

	return t, nil
}
