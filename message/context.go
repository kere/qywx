package message

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/util"
)

var (
	xmlContentType   = []string{"text/xml; charset=utf-8"}
	plainContentType = []string{"text/plain; charset=utf-8"}
)

// // ReplyMessageCall func
// type ReplyMessageCall func(msg MixMessage) (ICommonMessage, error)

// Context for message
type Context struct {
	Writer                          http.ResponseWriter
	Request                         *http.Request
	Nonce                           string
	Random                          []byte
	Timestamp                       string
	AppID, AppSecret, Token, AesKey string
	IsSafe                          bool

	MixMessage MixMessage
}

// NewContext func
func NewContext(w http.ResponseWriter, req *http.Request, appid, appsecret, token, aeskey string) Context {
	return Context{Writer: w, Request: req, AppID: appid, AppSecret: appsecret, AesKey: aeskey, Token: token}
}

// ParsePost 解析微信消息
func (c *Context) ParsePost() ([]byte, error) {
	req := c.Request

	var rawMsg []byte
	var err error
	if !c.IsSafe {
		c.IsSafe = req.FormValue("encrypt_type") == "aes"
	}

	if c.IsSafe {
		var eXML EncryptedXMLMsg
		err = xml.NewDecoder(req.Body).Decode(&eXML)
		if err != nil {
			return nil, err
		}

		sign := req.FormValue("msg_signature")
		c.Timestamp = req.FormValue("timestamp")
		c.Nonce = req.FormValue("nonce")

		devSign := util.Signature(c.Token, c.Timestamp, c.Nonce, eXML.EncryptedMsg)
		if devSign != sign {
			return nil, errors.New("sign failed")
		}

		c.Random, rawMsg, err = util.DecryptMsg(c.AppID, eXML.EncryptedMsg, c.AesKey)
		if err != nil {
			return nil, err
		}
	} else {
		rawMsg, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("从body中解析xml失败, err=%v", err)
		}
	}

	err = xml.Unmarshal(rawMsg, &c.MixMessage)
	return rawMsg, err
}

// Send 发送自定义回复
func (c Context) Send(msg IMessage) error {
	if msg.GetMsgType() == "" {
		return nil
	}

	if c.IsSafe {
		src, err := xml.Marshal(msg)
		if err != nil {
			return err
		}

		src, err = util.EncryptMsg(c.Random, src, c.AppID, c.AesKey)
		if err != nil {
			return err
		}

		ts := msg.GetCreateTime()
		if ts == 0 {
			ts = time.Now().Unix()
		}
		replyMsg := ResponseEncryptedXMLMsg{
			Nonce:     c.Nonce,
			Timestamp: ts,
		}

		encrypted := string(src)
		replyMsg.MsgSignature = util.Signature(c.Token, fmt.Sprint(replyMsg.Timestamp), replyMsg.Nonce, encrypted)
		replyMsg.EncryptedMsg = encrypted

		c.RenderXML(replyMsg)

	} else {
		c.RenderXML(msg)
	}

	log.App.Debug("send to user:", c.MixMessage.ToUserName)
	return nil
}

//Render render from bytes
func (c Context) Render(bytes []byte) {
	c.Writer.WriteHeader(200)
	_, err := c.Writer.Write(bytes)
	if err != nil {
		panic(err)
	}
}

//RenderString render from string
func (c Context) RenderString(str string) {
	writeContextType(c.Writer, plainContentType)
	c.Render([]byte(str))
}

//RenderXML render to xml
func (c Context) RenderXML(obj interface{}) {
	writeContextType(c.Writer, xmlContentType)
	bytes, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Render(bytes)
}

func writeContextType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
