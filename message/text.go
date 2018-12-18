package message

import (
	"time"
)

//Text 文本消息
type Text struct {
	CommonToken
	Content string `xml:"Content"`
}

//NewText 初始化文本消息
func NewText(content string) Text {
	text := Text{}
	text.Content = content
	text.MsgType = MsgTypeText

	return text
}

//NewEmpty 空文本消息
func NewEmpty() Text {
	text := Text{}
	text.MsgType = MsgTypeText
	return text
}

//NewReplyText 初始化文本消息
func NewReplyText(content string, mixmsg *MixMessage) Text {
	text := Text{}
	text.Content = content
	text.MsgType = MsgTypeText
	text.FromUserName = mixmsg.ToUserName
	text.ToUserName = mixmsg.GetFromUserName()
	text.CreateTime = time.Now().Unix()

	return text
}
