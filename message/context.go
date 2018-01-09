package message

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/util"
)

var (
	xmlContentType   = []string{"text/xml; charset=utf-8"}
	plainContentType = []string{"text/plain; charset=utf-8"}
)

// Context for message
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Agent   *corp.Agent
	Nonce   string

	MixMessage      *MixMessage
	EncryptedXMLMsg *EncryptedXMLMsg
}

// NewContext func
func NewContext(w http.ResponseWriter, req *http.Request, agent *corp.Agent) *Context {
	return &Context{Writer: w, Request: req, Agent: agent}
}

// AuthHTTPGet 微信消息接口验证
// 如果成功，返回 解密后的 字符串;不成功，返回空字符串
func (c *Context) AuthHTTPGet() string {
	// msg_signature: 企业微信加密签名，msg_signature结合了企业填写的token、请求中的timestamp、nonce参数、加密的消息体
	// timestamp: 时间戳
	// nonce: 随机数
	// echostr: 加密的随机字符串，以msg_encrypt格式提供。需要解密并返回echostr明文，解密后有random、msg_len、msg、$CorpID四个字段，其中msg即为echostr明文
	req := c.Request
	sign := req.FormValue("msg_signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")
	echostr := req.FormValue("echostr")
	// 企业计算签名：dev_msg_signature=sha1(sort(token、timestamp、nonce、msg_encrypt))。sort的含义是将参数值按照字母字典排序，然后从小到大拼接成一个字符串
	agent := c.Agent
	devSign := util.Signature(agent.MsgToken, timestamp, nonce, echostr)
	if devSign != sign {
		return ""
	}

	_, rawMsg, err := util.DecryptMsg(agent.Corp.ID, echostr, agent.MsgAesKey)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s", rawMsg)
}

// ParsePost 解析微信消息
func (c *Context) ParsePost() error {
	req := c.Request

	err := xml.NewDecoder(req.Body).Decode(&c.EncryptedXMLMsg)
	if err != nil {
		return err
	}
	agent := c.Agent

	sign := req.FormValue("msg_signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")
	devSign := util.Signature(agent.MsgToken, timestamp, nonce, c.EncryptedXMLMsg.EncryptedMsg)
	if devSign != sign {
		return errors.New("sign failed")
	}

	_, rawMsg, err := util.DecryptMsg(agent.Corp.ID, c.EncryptedXMLMsg.EncryptedMsg, agent.MsgAesKey)
	if err != nil {
		return err
	}

	c.Nonce = nonce
	err = xml.Unmarshal(rawMsg, &c.MixMessage)
	return err
}

// SendBy 发送自定义回复
func (c *Context) SendBy(f func(agent *corp.Agent, msg *MixMessage) (ICommonMessage, error)) error {
	if c.MixMessage == nil || c.MixMessage.ToUserName == "" {
		return nil
	}

	msg, err := f(c.Agent, c.MixMessage)
	if err != nil {
		return err
	}

	msg.SetCreateTime(time.Now().Unix())
	msg.SetFromUserName(c.Agent.Corp.ID)
	msg.SetToUserName(c.MixMessage.FromUserName)

	src, _ := xml.Marshal(msg)
	log.App.Debug("sendBy msg:", string(src))

	rmsg, err := c.buildResponseEncryptedMsg(msg)
	if err != nil {
		return err
	}

	src, _ = xml.Marshal(rmsg)
	log.App.Debug("sendBy encrypt:", string(src))

	WriteXML(c.Writer, rmsg)
	return nil
}

// BuildResponseEncryptedMsg 返回创建好的 ResponseEncryptedXMLMsg
func (c *Context) buildResponseEncryptedMsg(v interface{}) (replyMsg ResponseEncryptedXMLMsg, err error) {
	unix := time.Now().Unix()
	replyMsg.Nonce = c.Nonce
	replyMsg.Timestamp = unix

	agent := c.Agent

	src, err := xml.Marshal(v)
	if err != nil {
		return replyMsg, err
	}

	src, err = util.EncryptMsg([]byte(c.Nonce), src, agent.Corp.ID, agent.MsgAesKey)
	if err != nil {
		return replyMsg, err
	}

	encrypted := string(src)
	replyMsg.MsgSignature = util.Signature(agent.MsgToken, fmt.Sprint(replyMsg.Timestamp), replyMsg.Nonce, encrypted)

	replyMsg.EncryptedMsg = encrypted
	return replyMsg, err
}

// WriteXML 写入xml信息
func WriteXML(w http.ResponseWriter, obj interface{}) {
	if obj == nil {
		return
	}

	util.WriteContextType(w, xmlContentType)
	bytes, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}

// WriteString 写入string信息
func WriteString(w http.ResponseWriter, txt string) {
	util.WriteContextType(w, plainContentType)
	w.WriteHeader(200)
	_, err := w.Write([]byte(txt))
	if err != nil {
		panic(err)
	}
}

// // CDATA wrap
// func CDATA(s string) string {
// 	return s
// 	// return "<![CDATA[" + s + "]]>"
// }
