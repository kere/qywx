package corp

import (
	"fmt"

	"github.com/kere/qywx/client"
)

const (
	agentInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/agent/get?access_token=%s&agentid=%d"
)

// AgentInfo class
type AgentInfo struct {
	AgentID        int      `json:"agentid"`
	Name           string   `json:"name"`
	LogURL         string   `json:"square_log_url"`
	Desc           string   `json:"description"`
	AllowUserIDs   []string `json:"allow_userid"`
	AllowPartys    []int    `json:"allow_partyid"`
	AllowTags      []int    `json:"allow_tags"`
	Close          int      `json:"close"`
	RedirectDomain string
	ReportLocation int `json:"report_location_flag"`
	//是否上报用户进入应用事件。0：不接收；1：接收
	IsReportenter int    `json:"isreportenter"`
	HomeURL       string `json:"home_url"`
}

// WXAgentInfo 获得微信应用信息,企业仅可获取当前凭证对应的应用
func WXAgentInfo(agentID int, token string) (a AgentInfo, err error) {
	dat, err := client.Get(fmt.Sprintf(agentInfoURL, token, agentID))
	if err != nil {
		return a, err
	}

	a.AgentID = dat.Int("agentid")
	a.Name = dat.String("name")
	a.LogURL = dat.String("square_logo_url")
	a.Desc = dat.String("description")

	// "allow_userinfos":{
	//   "user":[]
	if dat.IsSet("allow_userinfos") {
		obj := dat["allow_userinfos"].(map[string]interface{})
		if users, isok := obj["user"]; isok {
			usrs := users.([]interface{})
			for _, u := range usrs {
				v := (u.(map[string]interface{}))["userid"].(string)
				a.AllowUserIDs = append(a.AllowUserIDs, v)
			}
		}
	}

	if dat.IsSet("allow_partys") {
		arr := (dat["allow_partys"].(map[string]interface{}))["partyid"].([]interface{})

		a.AllowPartys = make([]int, len(arr))
		for i := range arr {
			a.AllowPartys[i] = int(arr[i].(float64))
		}
	}

	if dat.IsSet("allow_tags") {
		arr := (dat["allow_tags"].(map[string]interface{}))["tagid"].([]interface{})
		a.AllowTags = make([]int, len(arr))
		for i := range arr {
			a.AllowTags[i] = int(arr[i].(float64))
		}
	}

	a.Close = dat.Int("close")
	a.RedirectDomain = dat.String("redirect_domain")
	a.ReportLocation = dat.Int("report_location_flag")
	a.IsReportenter = dat.Int("isreportenter")
	a.HomeURL = dat.String("home_url")

	return a, err
}
