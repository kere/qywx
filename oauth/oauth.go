package oauth

import (
	"fmt"
	"net/http"

	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
)

// OAuth class
type OAuth struct {
	Appid        string
	ResponseType string
	Scope        string

	AgentID string
	State   string
}

// NewOAuth func
func NewOAuth(appid, agentID string) *OAuth {
	return &OAuth{Appid: appid, AgentID: agentID, ResponseType: "code", Scope: "snsapi_base"}
}

// OpenAndRedirect oauth
func (o *OAuth) OpenAndRedirect(rw http.ResponseWriter, req *http.Request, redirectURL string) {
	url := fmt.Sprintf(openURL, o.Appid, redirectURL, o.ResponseType, o.Scope, o.AgentID, o.State)

	http.Redirect(rw, req, url, http.StatusSeeOther)
}

// GetUserInfo 根据code获取成员信息
// 请求方式：GET（HTTPS）
func (o *OAuth) GetUserInfo(accessToken, code string) (*UserInfo, error) {
	url := fmt.Sprintf(userInfoURL, accessToken, code)
	dat, err := client.GetMapData(url)
	if err != nil {
		return nil, err
	}

	if dat.IsSet("OpenId") {
		// 非企业用户
		return &UserInfo{ID: dat.String("OpenId"), DevicedID: dat.String("DeviceId"), IsOpenUser: true}, nil
	}

	// 企业用户
	usr := &UserInfo{ID: dat.String("UserId"), DevicedID: dat.String("DeviceId"), Ticket: dat.String("user_ticket")}
	usr.SetTicketExpires(dat.Int("expires_in"))

	return usr, nil
}

// GetUserDetail 使用user_ticket获取成员详情
// 请求方式：POST（HTTPS）
// Post: "user_ticket": "USER_TICKET"
func (o *OAuth) GetUserDetail(accessToken, ticket string) (*UserDetail, error) {
	url := fmt.Sprintf(userDetailURL, accessToken)
	dat, err := client.PostMapData(url, util.MapData{"user_ticket": ticket})
	if err != nil {
		return nil, err
	}

	usr := &UserDetail{Name: dat.String("name"), ID: dat.String("userid"), Position: dat.String("position"), Mobile: dat.String("mobile"), Gender: dat.Int("gender"), Email: dat.String("email"), Avatar: dat.String("avatar"), Department: dat.Ints("department")}
	return usr, nil
}
