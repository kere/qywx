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
func SetUser(corpID string, usrDetail users.UserDetail, expires int) {
	if usrDetail.ID == "" {
		return
	}

	src, _ := json.Marshal(usrDetail)
	v := md5.Sum([]byte(corpID + usrDetail.ID))
	cache.Set(fmt.Sprintf("%s-%x", usrKey, v), string(src), expires)
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
