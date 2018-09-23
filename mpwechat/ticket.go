package mpwechat

import (
	"fmt"
	"strconv"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/client"
)

const (
	cacheJSTicketKeyPrefix = "jsticket:"
)

// WxJSTicket 获得微信JS访问ticket
func WxJSTicket(appid, appsecret string) (string, error) {
	key := cacheJSTicketKeyPrefix + appid
	val, err := cache.GetString(key)
	if err == nil {
		return val, err
	}
	if val != "" {
		return val, nil
	}

	//从微信服务器获取
	var expires int64
	val, expires, err = loadJSticket(appid, appsecret)
	if err != nil {
		return val, err
	}

	expires -= 600
	if expires < 0 {
		expires = 0
	}

	cache.Set(key, val, int(expires))
	return val, nil
}

//loadAccessToken 获取token
func loadJSticket(appid, appsecret string) (string, int64, error) {
	token, err := WxAccessToken(appid, appsecret)
	if err != nil {
		return "", 0, err
	}
	uri := fmt.Sprintf(WxURL("GetJSTicket"), token)
	dat, err := client.Get(uri, nil)
	if err != nil {
		return "", 0, err
	}
	expires, err := strconv.ParseInt(dat.String("expires_in"), 10, 64)
	if err != nil {
		return "", 0, err
	}

	return dat.String("ticket"), expires, nil
}
