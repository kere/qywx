package message

import (
	"encoding/json"

	"github.com/kere/gno/db"
)

const (
	// ChangeTypeCreateUser 新建用户
	ChangeTypeCreateUser = "create_user"
	// ChangeTypeUpdateUser 更新用户
	ChangeTypeUpdateUser = "update_user"
	// ChangeTypeDeleteUser 删除用户
	ChangeTypeDeleteUser = "delete_user"
	// ChangeTypeUpdateTag update tag
	ChangeTypeUpdateTag = "update_tag"
	// ChangeTypeDeleteTag delete tag
	ChangeTypeDeleteTag = "delete_tag"
	// ChangeTypeUpdateParty update party
	ChangeTypeUpdateParty = "update_party"
	// ChangeTypeDeleteParty delete party
	ChangeTypeDeleteParty = "delete_party"

	//EventChangeContact 订阅
	EventChangeContact = "change_contact"
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

// ReplyContactEventCall func
type ReplyContactEventCall func(e ContactEvent) error

// // CommonToken 消息中通用的结构
// type CommonToken struct {
// 	XMLName      xml.Name `xml:"xml"`
// 	ToUserName   string   `xml:"ToUserName"`
// 	FromUserName string   `xml:"FromUserName"`
// 	CreateTime   int64    `xml:"CreateTime"`
// 	MsgType      string   `xml:"MsgType"`
// 	Event        string   `xml:"Event"`
// }

//ContactEvent 通讯录事件
type ContactEvent struct {
	CommonToken

	ChangeType string `xml:"ChangeType"`

	UserID    string `xml:"UserID"`
	NewUserID string `xml:"NewUserID"`

	Name          string `xml:"Name"`
	Department    string `xml:"Department"`
	Mobile        string `xml:"Mobile"`
	Position      string `xml:"Position"`
	Gender        string `xml:"Gender"`
	Email         string `xml:"Email"`
	Status        string `xml:"Status"`
	Avatar        string `xml:"Avatar"`
	EnglishName   string `xml:"EnglishName"`
	IsLeader      string `xml:"IsLeader"`
	Telephone     string `xml:"Telephone"`
	TagID         string `xml:"TagId"`
	AddUserItems  string `xml:"AddUserItems"`
	DelUserItems  string `xml:"DelUserItems"`
	AddPartyItems string `xml:"AddPartyItems"`
	DelPartyItems string `xml:"DelPartyItems"`
}

// ToDataRow to db.DataRow
func (c ContactEvent) ToDataRow() db.DataRow {
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
		var v []int
		src := "[" + c.Department + "]"
		err := json.Unmarshal([]byte(src), &v)
		if err == nil {
			row["department"] = v
		}
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
