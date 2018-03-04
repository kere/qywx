package oa

import (
	"encoding/json"
	"fmt"

	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
)

const (
	checkinURL = "https://qyapi.weixin.qq.com/cgi-bin/checkin/getcheckindata?access_token=%s"
)

// CheckinOA class
type CheckinOA struct {
	UserID         string `json:"userid"`
	GroupName      string `json:"groupname"`
	CheckinType    string `json:"checkin_type"`
	CheckinTime    int64  `json:"checkin_time"`
	ExceptionType  string `json:"exception_type"`
	LocationTitle  string `json:"location_title"`
	LocationDetail string `json:"location_detail"`
	WifiName       string `json:"wifiname"`
	Notes          string `json:"notes"`
	WifiMAC        string `json:"wifimac"`
	MediaID        string `json:"mediaid"`
}

// WxCheckin get list
// 打卡类型 itype 1：上下班打卡；2：外出打卡；3：全部打卡
func WxCheckin(useridlist []string, start, end int64, itype int, token string) ([]CheckinOA, error) {
	items := make([]CheckinOA, 0)
	args := util.MapData{
		"opencheckindatatype": itype,
		"starttime":           start,
		"endtime":             end,
		"useridlist":          useridlist,
	}

	dat, err := client.PostJSON(fmt.Sprintf(checkinURL, token), args)
	if err != nil {
		return items, err
	}

	src, err := json.Marshal(dat["checkindata"])
	if err != nil {
		return items, err
	}

	err = json.Unmarshal(src, &items)

	return items, err
}
