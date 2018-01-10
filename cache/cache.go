package cache

import (
	"encoding/json"
	"sync"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/conf"
	"github.com/kere/qywx/users"
)

const (
	usrKey = "wxu-"
)

var (
	usrMutex = new(sync.RWMutex)
)

// SetUser detail
func SetUser(corpID string, usr users.UserDetail, expires int) error {
	if usr.ID == "" {
		return nil
	}

	src, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	return cache.Set(Key(corpID, usr.ID), string(src), expires)
}

// Key cached
func Key(corpID, uid string) string {
	return usrKey + corpID + "-" + uid
}

// GetUser detail
func GetUser(corpID, uid string) (usr users.UserDetail) {
	src, _ := cache.Get(Key(corpID, uid))
	if len(src) == 0 {
		return usr
	}

	json.Unmarshal([]byte(src), &usr)

	return usr
}

// Init redis instance
func Init(config conf.Conf) {
	cache.Init(config)
}
