package users

import (
	"errors"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/util"
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
	t.Init(t, 0)
	return t
}

// GetUserDetail func
func GetUserDetail(corpIndex int, userid string) UserDetail {
	v := cachedusermap.Get(corpIndex, userid)
	if v == nil {
		return UserDetail{}
	}
	return v.(UserDetail)
}

// ReleaseUserCache cache
func ReleaseUserCache(corpIndex int, userid string) {
	cachedusermap.Release(corpIndex, userid)
}

// Build func
func (t *cachedUserMap) Build(args ...interface{}) (interface{}, error) {
	corpIndex := args[0].(int)
	userid := args[1].(string)

	cp := corp.Get(corpIndex)
	if cp == nil {
		return nil, errors.New("corp not found in departCached")
	}

	token, err := cp.GetContactToken()
	if err != nil {
		return nil, err
	}

	dat, err := WxUser(userid, token)
	if err != nil {
		return nil, err
	}

	t.SetExpires(util.Expires())
	return dat, nil
}
