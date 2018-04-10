package oauth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
	"github.com/kere/qywx/users"
)

// OAuth class
type OAuth struct {
	Appid        string
	ResponseType string
	Scope        string

	AgentID int
	State   string
}

// NewOAuth func
func NewOAuth(appid string, agentID int) *OAuth {
	return &OAuth{Appid: appid, AgentID: agentID, ResponseType: "code", Scope: "snsapi_base"}
}

// OpenAndRedirect oauth
func (o *OAuth) OpenAndRedirect(rw http.ResponseWriter, req *http.Request, redirectURL string) {
	url := fmt.Sprintf(openURL, o.Appid, redirectURL, o.ResponseType, o.Scope, string(o.AgentID), o.State)

	http.Redirect(rw, req, url, http.StatusSeeOther)
}

// GetUserInfo 根据code获取成员信息
// 请求方式：GET（HTTPS）
func (o *OAuth) GetUserInfo(accessToken, code string) (usr UserInfo, err error) {
	url := fmt.Sprintf(userInfoURL, accessToken, code)
	dat, err := client.Get(url)
	if err != nil {
		return usr, err
	}

	log.App.Debug("userinfo", dat)

	if dat.IsSet("OpenId") {
		// 非企业用户
		return UserInfo{
			UserID:     dat.String("OpenId"),
			DevicedID:  dat.String("DeviceId"),
			IsOpenUser: true}, nil
	}

	// 企业用户
	usr = UserInfo{
		UserID:    dat.String("UserId"),
		DevicedID: dat.String("DeviceId")}

	if dat.IsSet("user_ticket") {
		usr.Ticket = dat.String("user_ticket")
	}

	if dat.IsSet("expires_in") {
		usr.SetTicketExpires(dat.Int("expires_in"))
	}

	return usr, nil
}

// GetUserDetail 使用user_ticket获取成员详情
// 请求方式：POST（HTTPS）
// Post: "user_ticket": "USER_TICKET"
func (o *OAuth) GetUserDetail(accessToken, ticket string) (usr users.UserDetail, err error) {
	if ticket == "" {
		return usr, errors.New("user ticket is empty")
	}

	url := fmt.Sprintf(userDetailURL, accessToken)
	dat, err := client.PostJSON(url, util.MapData{"user_ticket": ticket})
	if err != nil {
		return usr, err
	}

	log.App.Debug("userdetail", dat)

	usr = users.UserDetail{Name: dat.String("name"), UserID: dat.String("userid"), Position: dat.String("position"), Mobile: dat.String("mobile"), Gender: dat.String("gender"), Email: dat.String("email"), Avatar: dat.String("avatar"), Department: dat.Ints("department")}
	return usr, nil
}
