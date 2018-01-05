package oauth

import (
	"errors"
	"time"

	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/oauth"
)

const (
	openURL = "    https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=%s&scope=%s&agentid=%s&state=%s#wechat_redirect"

	userInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s"

	userDetailURL = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserdetail?access_token=%s"
)

// UserInfo class
type UserInfo struct {
	IsOpenUser            bool
	ID, DevicedID, Ticket string
	TicketExpires         int
	TicketExpiresAt       time.Time
}

// SetTicketExpires build expires
func (u *UserInfo) SetTicketExpires(n int) {
	u.TicketExpires = n
	u.TicketExpiresAt = time.Now().Add(time.Duration(n) * time.Second)
}

// UserDetail class
type UserDetail struct {
	ID         string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"depat"`
	Position   string `json:"position"` // 职位
	Mobile     string
	Gender     int `json:"gender"` // 性别
	Email      string
	Avatar     string `json:"avatar"`
}

func FetchUser(agentName, code string) (userInfo *oauth.UserInfo, userDetail *oauth.UserDetail, err error) {
	if code == "" {
		return nil, nil, errors.New("user code is empty")
	}

	a := corp.Corp.GetAgent(agentName)
	token, err := a.GetToken()
	if err != nil {
		return nil, nil, err
	}

	oa := oauth.NewOAuth(corp.Corp.ID, a.ID)
	oa.State = ""

	userInfo, err = oa.GetUserInfo(token.Value, code)
	if err != nil {
		return nil, nil, err
	}

	if userInfo.IsOpenUser {
		return userInfo, nil, nil
	}

	userDetail, err = oa.GetUserDetail(token.Value, userInfo.Ticket)
	return userInfo, userDetail, err
}
