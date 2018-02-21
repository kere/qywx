package cached

import (
	"errors"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/users"
)

// cachedusermap
var (
	// cachedusermap a
	cachedusermap = newCachedUserMap()
)

// cachedUserMap class
type cachedUserMap struct {
	cache.Map
}

func newCachedUserMap() *cachedUserMap {
	t := &cachedUserMap{}
	t.Init(t)
	return t
}

// GetUserDetail func
func GetUserDetail(corpIndex int, userid string) users.UserDetail {
	v := cachedusermap.Get(corpIndex, userid)
	if v == nil {
		return users.UserDetail{}
	}
	return v.(users.UserDetail)
}

// ReleaseUserCache cache
func ReleaseUserCache(corpIndex int, userid string) {
	cachedusermap.Release(corpIndex, userid)
}

// Build func
func (t *cachedUserMap) Build(args ...interface{}) (interface{}, int, error) {
	corpIndex := args[0].(int)
	userid := args[1].(string)

	cp := corp.Get(corpIndex)
	if cp == nil {
		return nil, 0, errors.New("corp not found in departCached")
	}

	token, err := cp.GetContactToken()
	if err != nil {
		return nil, 0, err
	}

	dat, err := users.WxUser(userid, token)
	if err != nil {
		return nil, 0, err
	}

	return dat, Expires(), nil
}
