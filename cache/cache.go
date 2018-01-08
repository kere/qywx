package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
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

	src, _ := json.Marshal(usr)
	v := md5.Sum([]byte(corpID + usr.ID))

	err := cache.Set(fmt.Sprintf("%s-%x", usrKey, v), string(src), expires)
	return err
}

// GetUser detail
func GetUser(corpID, uid string) (usr users.UserDetail) {
	v := md5.Sum([]byte(corpID + uid))
	src, _ := cache.Get(fmt.Sprintf("%s-%x", usrKey, v))
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
