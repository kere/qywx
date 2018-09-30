package corp

import (
	"errors"
	"fmt"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/util"
)

type tokenContactCached struct {
	cache.Map
}

func newTokenContactCached() *tokenContactCached {
	t := &tokenContactCached{}
	t.Init(t, 0)
	return t
}

// CheckValue 检查缓存值是否正确，如果正确缓存
func (t *tokenContactCached) CheckValue(v interface{}) bool {
	return v.(string) != ""
}

// Build func
func (t *tokenContactCached) Build(args ...interface{}) (interface{}, error) {
	corpName := args[0].(string)

	cp, err := GetByName(corpName)
	if err != nil {
		return nil, errors.New("corp not found in tokenContactCached")
	}

	// 获取 access_token
	// 请求方式：GET（HTTPS）
	dat, err := util.AjaxGet(fmt.Sprintf(tokenURL, cp.Corpid, cp.ContactsSecret), nil)
	if err != nil {
		return "", err
	}

	t.SetExpires(dat.Int("expires_in"))
	// v = newToken(dat.String("access_token"), dat.Int("expires_in"))
	return dat.String("access_token"), nil
}

// map[agentID] *tokenContactCached
var tkContactCached = newTokenContactCached()

// GetContactToken get cached token
func (c *Corporation) GetContactToken() (string, error) {
	token := tkContactCached.Get(c.Name)
	if token == nil {
		return "", errors.New("get cached contact ticket is nil")
	}

	return token.(string), nil
}
