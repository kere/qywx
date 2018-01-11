package oauth

import (
	"errors"
	"time"

	"github.com/kere/qywx/users"
)

const (
	openURL = "    https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=%s&scope=%s&agentid=%s&state=%s#wechat_redirect"

	userInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s"

	userDetailURL = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserdetail?access_token=%s"
)

// UserInfo class
type UserInfo struct {
	IsOpenUser      bool   `json:"is_openuser"`
	ID              string `json:"userid"`
	DevicedID       string `json:"deviced_id"`
	Ticket          string `json:"ticket"`
	TicketExpires   int    `json:"expires"`
	TicketExpiresAt time.Time
}

// SetTicketExpires build expires
func (u *UserInfo) SetTicketExpires(n int) {
	u.TicketExpires = n
	u.TicketExpiresAt = time.Now().Add(time.Duration(n) * time.Second)
}

// FetchUser get userInfo & userDetail
func FetchUser(corpID string, agentID int, code, token, scope string) (userInfo UserInfo, userDetail users.UserDetail, err error) {
	if code == "" {
		return userInfo, userDetail, errors.New("user code is empty")
	}

	oa := NewOAuth(corpID, agentID)
	if scope != "" {
		oa.Scope = scope
	}
	oa.State = ""

	userInfo, err = oa.GetUserInfo(token, code)
	if err != nil {
		return userInfo, userDetail, err
	}

	if userInfo.IsOpenUser {
		return userInfo, userDetail, errors.New("openuser")
	}

	if userInfo.Ticket != "" {
		userDetail, err = oa.GetUserDetail(token, userInfo.Ticket)
	}

	return userInfo, userDetail, err
}
