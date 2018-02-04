package event

import (
	"encoding/xml"
	"strings"

	"github.com/kere/gno/db"
)

const (
	//EventSubscribe 订阅
	EventSubscribe = "subscribe"
	//EventUnsubscribe 取消订阅
	EventUnsubscribe = "unsubscribe"
	//EventScan 用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	EventScan = "SCAN"
	//EventLocation 上报地理位置事件
	EventLocation = "LOCATION"
	//EventClick 点击菜单拉取消息时的事件推送
	EventClick = "CLICK"
	//EventView 点击菜单跳转链接时的事件推送
	EventView = "VIEW"
	//EventScancodePush 扫码推事件的事件推送
	EventScancodePush = "scancode_push"
	//EventScancodeWaitmsg 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventScancodeWaitmsg = "scancode_waitmsg"
	//EventPicSysphoto 弹出系统拍照发图的事件推送
	EventPicSysphoto = "pic_sysphoto"
	//EventPicPhotoOrAlbum 弹出拍照或者相册发图的事件推送
	EventPicPhotoOrAlbum = "pic_photo_or_album"
	//EventPicWeixin 弹出微信相册发图器的事件推送
	EventPicWeixin = "pic_weixin"
	//EventLocationSelect 弹出地理位置选择器的事件推送
	EventLocationSelect = "location_select"
)

// CommonToken 消息中通用的结构
type CommonToken struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
}

//ContactsEvent 通讯录事件
type ContactsEvent struct {
	CommonToken

	ChangeType string `xml:"ChangeType"`

	UserID    string `xml:"UserID"`
	NewUserID string `xml:"NewUserID"`

	Name        string `xml:"Name"`
	Department  string `xml:"Department"`
	Mobile      string `xml:"Mobile"`
	Position    string `xml:"Position"`
	Gender      string `xml:"Gender"`
	Email       string `xml:"Email"`
	Status      string `xml:"Status"`
	Avatar      string `xml:"Avatar"`
	EnglishName string `xml:"EnglishName"`
	IsLeader    string `xml:"IsLeader"`
	Telephone   string `xml:"Telephone"`
}

// ToDataRow to db.DataRow
func (c ContactsEvent) ToDataRow() db.DataRow {
	if c.UserID == "" {
		return nil
	}
	row := db.DataRow{"userid": c.UserID}

	if c.Mobile != "" {
		row["mobile"] = c.Mobile
	}
	if c.NewUserID != "" {
		row["userid_old"] = c.UserID
		row["userid"] = c.NewUserID
	}
	if c.Email != "" {
		row["email"] = c.Email
	}
	if len(c.Department) > 0 {
		row["department"] = strings.Split(c.Department, ",")
	}
	if c.Name != "" {
		row["name"] = c.Name
	}
	if c.Avatar != "" {
		row["avatar"] = c.Avatar
	}
	if c.Position != "" {
		row["position"] = c.Position
	}
	if c.Gender != "" {
		row["gender"] = c.Gender
	}
	if c.Status != "" && c.Status != "0" {
		row["status"] = c.Status
	}

	if c.IsLeader != "" && c.IsLeader != "0" {
		row["isleader"] = c.IsLeader
	}
	return row
}
