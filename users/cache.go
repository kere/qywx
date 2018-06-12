package users

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/conf"
	"github.com/kere/qywx/corp"
)

const (
	usrKey = "wxu-"
)

var (
	usrMutex = new(sync.RWMutex)
)

// CacheSetUser detail
func CacheSetUser(agent *corp.Agent, usr UserDetail, expires int) error {
	if usr.UserID == "" {
		return nil
	}

	src, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	return cache.Set(CacheKey(agent, usr.UserID), string(src), expires)
}

// CacheKey cached
func CacheKey(agent *corp.Agent, uid string) string {
	return usrKey + fmt.Sprint(agent.Corp.ID) + "-" + fmt.Sprint(agent.ID) + "-" + uid
}

// CacheGetUser detail
func CacheGetUser(agent *corp.Agent, uid string) (usr UserDetail) {
	src, _ := cache.Get(CacheKey(agent, uid))
	if len(src) == 0 {
		return usr
	}

	json.Unmarshal([]byte(src), &usr)

	return usr
}

// CacheDelUser clear
func CacheDelUser(agent *corp.Agent, uid string) error {
	return cache.Delete(CacheKey(agent, uid))
}

// Init redis instance
func Init(config conf.Conf) {
	cache.Init(config)
}
