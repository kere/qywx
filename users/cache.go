package users

import (
	"encoding/json"
	"sync"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/conf"
)

const (
	usrKey = "wxu-"
)

var (
	usrMutex = new(sync.RWMutex)
)

// CacheSetUser detail
func CacheSetUser(corpID string, usr UserDetail, expires int) error {
	if usr.UserID == "" {
		return nil
	}

	src, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	return cache.Set(CacheKey(corpID, usr.UserID), string(src), expires)
}

// CacheKey cached
func CacheKey(corpID, uid string) string {
	return usrKey + corpID + "-" + uid
}

// CacheGetUser detail
func CacheGetUser(corpID, uid string) (usr UserDetail) {
	src, _ := cache.Get(CacheKey(corpID, uid))
	if len(src) == 0 {
		return usr
	}

	json.Unmarshal([]byte(src), &usr)

	return usr
}

// CacheDelUser detail
func CacheDelUser(corpID, uid string) error {
	return cache.Delete(CacheKey(corpID, uid))
}

// Init redis instance
func Init(config conf.Conf) {
	cache.Init(config)
}
