package oa

import (
	"encoding/json"
	"fmt"
	"time"

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

// ToSimple f
func (c CheckinOA) ToSimple() *CheckinSimpleOA {
	return &CheckinSimpleOA{
		UserID:      c.UserID,
		GroupName:   c.GroupName,
		CheckinType: c.CheckinType,
		CheckinTime: c.CheckinTime,
	}
}

// CheckinSorted data
type CheckinSorted []CheckinOA

func (a CheckinSorted) Len() int {
	return len(a)
}
func (a CheckinSorted) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a CheckinSorted) Less(i, j int) bool {
	return a[i].CheckinTime > a[j].CheckinTime
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

// CheckinSimpleOA class
type CheckinSimpleOA struct {
	UserID      string `json:"userid"`
	GroupName   string `json:"groupname"`
	CheckinType string `json:"checkin_type"`
	CheckinTime int64  `json:"checkin_time"`
}

// WorkData class
type WorkData struct {
	Start *CheckinSimpleOA `json:"start"`
	Off   *CheckinSimpleOA `json:"off"`
	Hours float64          `json:"hours"`
}

// Build f
func (w *WorkData) Build() {
	a := time.Unix(w.Start.CheckinTime, 0)
	b := time.Unix(w.Off.CheckinTime, 0)
	d := b.Sub(a)
	w.Hours = float64(d / time.Hour)
}

// GroupCheckin class
type GroupCheckin struct {
	UserID   string             `json:"userid"`
	UserName string             `json:"user_name"`
	Items    []*CheckinSimpleOA `json:"-"`
	Works    []*WorkData
	Memo     string `json:"memo"`
}

const (
	gowork  = "上班打卡"
	offwork = "下班打卡"
)

// GroupCheckinData f
func GroupCheckinData(records []CheckinOA) []GroupCheckin {
	var tmp = make(map[string][]*CheckinSimpleOA)
	for _, r := range records {
		if r.ExceptionType != "" {
			continue
		}
		if _, isok := tmp[r.UserID]; !isok {
			tmp[r.UserID] = []*CheckinSimpleOA{}
		}
		tmp[r.UserID] = append(tmp[r.UserID], r.ToSimple())
	}

	var result = []GroupCheckin{}
	for k, items := range tmp {
		g := GroupCheckin{}
		g.UserID = k
		g.Items = items

		work := &WorkData{}
		for _, v := range items {
			if v.CheckinType == gowork && work.Start == nil && work.Off == nil {
				work.Start = v
				continue
			}

			if v.CheckinType == offwork && work.Start != nil && work.Off == nil {
				if !isSameDay(work.Start, v) {
					work = &WorkData{}
					continue
				}
				work.Off = v
				work.Build()
				g.Works = append(g.Works, work)
				work = &WorkData{}
			}
		}

		result = append(result, g)
	}

	return result
}

func isSameDay(a, b *CheckinSimpleOA) bool {
	ta := time.Unix(a.CheckinTime, 0)
	tb := time.Unix(b.CheckinTime, 0)
	return ta.Format("2006-01-02") == tb.Format("2006-01-02")
}
