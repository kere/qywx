package mpwechat

import (
	"fmt"
	"strconv"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/qywx/util"
)

const (
	cacheTokenKeyPrefix = "acctoken:"
)

// WxAccessToken 获得微信访问Token
func WxAccessToken(appid, appsecret string) (string, error) {
	// 获取 access_token
	// 请求方式：GET（HTTPS）
	accessTokenCacheKey := cacheTokenKeyPrefix + appid
	val, err := cache.GetString(accessTokenCacheKey)
	if err == nil {
		return val, err
	}
	if val != "" {
		return val, nil
	}

	//从微信服务器获取
	var expires int64
	val, expires, err = loadAccessToken(appid, appsecret)
	if err != nil {
		return val, err
	}

	expires -= 600
	if expires < 0 {
		expires = 0
	}

	cache.Set(accessTokenCacheKey, val, int(expires))
	return val, nil
}

//loadAccessToken 获取token
func loadAccessToken(appid, appsecret string) (string, int64, error) {
	uri := fmt.Sprintf(WxURL("Token"), appid, appsecret)

	dat, err := util.AjaxGet(uri, nil)
	if err != nil {
		return "", 0, err
	}
	expires, err := strconv.ParseInt(dat.String("expires_in"), 10, 64)
	if err != nil {
		return "", 0, err
	}

	return dat.String("access_token"), expires, nil
}
