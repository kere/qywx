package mpwechat

import (
	"fmt"

	"github.com/kere/gno/libs/util"

	wxutil "github.com/kere/qywx/util"
)

// WxUserInfo 获得微信用户信息
func WxUserInfo(appid, appsecret, openid string) (util.MapData, error) {
	acctoken, err := WxAccessToken(appid, appsecret)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf(WxURL("UserInfo"), acctoken, openid)

	return wxutil.AjaxGet(uri, nil)
}
