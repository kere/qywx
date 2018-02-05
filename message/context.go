package message

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/util"
)

// ReplyMessageCall func
type ReplyMessageCall func(msg MixMessage) (ICommonMessage, error)

// Context for message
type Context struct {
	Writer                http.ResponseWriter
	Request               *http.Request
	Nonce                 string
	CorpID, Token, AesKey string

	MixMessage MixMessage
}

// NewContext func
func NewContext(w http.ResponseWriter, req *http.Request, corpID, token, aeskey string) Context {
	return Context{Writer: w, Request: req, CorpID: corpID, AesKey: aeskey, Token: token}
}

// ParsePost 解析微信消息
func (c *Context) ParsePost() ([]byte, error) {
	req := c.Request

	var eXML EncryptedXMLMsg
	err := xml.NewDecoder(req.Body).Decode(&eXML)
	if err != nil {
		return nil, err
	}

	sign := req.FormValue("msg_signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")

	devSign := util.Signature(c.Token, timestamp, nonce, eXML.EncryptedMsg)
	if devSign != sign {
		return nil, errors.New("sign failed")
	}

	_, rawMsg, err := util.DecryptMsg(c.CorpID, eXML.EncryptedMsg, c.AesKey)
	if err != nil {
		return nil, err
	}
	c.Nonce = nonce

	log.App.Debug(string(rawMsg))
	err = xml.Unmarshal(rawMsg, &c.MixMessage)
	return rawMsg, err
}

// SendBy 发送自定义回复
func (c *Context) SendBy(f ReplyMessageCall) error {
	if c.MixMessage.ToUserName == "" {
		return nil
	}

	msg, err := f(c.MixMessage)
	if err != nil {
		return err
	}

	msg.SetCreateTime(time.Now().Unix())
	msg.SetFromUserName(c.CorpID)
	msg.SetToUserName(c.MixMessage.FromUserName)

	rmsg, err := c.buildResponseEncryptedMsg(msg)
	if err != nil {
		return err
	}

	log.App.Debug("sendmsg:", c.MixMessage.ToUserName)
	util.WriteXML(c.Writer, rmsg)
	return nil
}

// BuildResponseEncryptedMsg 返回创建好的 ResponseEncryptedXMLMsg
func (c *Context) buildResponseEncryptedMsg(v interface{}) (replyMsg ResponseEncryptedXMLMsg, err error) {
	unix := time.Now().Unix()
	replyMsg.Nonce = c.Nonce
	replyMsg.Timestamp = unix

	src, err := xml.Marshal(v)
	if err != nil {
		return replyMsg, err
	}

	src, err = util.EncryptMsg([]byte(c.Nonce), src, c.CorpID, c.AesKey)
	if err != nil {
		return replyMsg, err
	}

	encrypted := string(src)
	replyMsg.MsgSignature = util.Signature(c.Token, fmt.Sprint(replyMsg.Timestamp), replyMsg.Nonce, encrypted)

	replyMsg.EncryptedMsg = encrypted
	return replyMsg, err
}
