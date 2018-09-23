package mpwechat

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/kere/gno/libs/util"
	wxutil "github.com/kere/qywx/util"
)

// WxJSConfig 微信JS验证数据
func WxJSConfig(appid, appsecret, urlstr string) (util.MapData, error) {
	ticket, err := WxJSTicket(appid, appsecret)
	if err != nil {
		return nil, err
	}

	// signature
	nonce := wxutil.RandomStr(16)
	timestamp := fmt.Sprint(time.Now().Unix())

	arr := []string{
		"jsapi_ticket=" + ticket,
		"noncestr=" + nonce,
		"timestamp=" + timestamp,
		"url=" + urlstr,
	}

	h := sha1.New()
	io.WriteString(h, strings.Join(arr, "&"))
	sign := fmt.Sprintf("%x", h.Sum(nil))

	data := util.MapData{
		"signature": sign,
		"nonce":     nonce,
		"appid":     appid,
		"timestamp": timestamp,
		"url":       urlstr}

	return data, nil
}
