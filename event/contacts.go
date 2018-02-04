package event

import (
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/message"
	"github.com/kere/qywx/util"
)

// ContactsContext for message
type ContactsContext struct {
	Writer                http.ResponseWriter
	Request               *http.Request
	CorpID, Token, AesKey string
	Nonce                 string

	Event ContactsEvent
}

// NewContactsContext func
func NewContactsContext(w http.ResponseWriter, req *http.Request, corpID, token, aeskey string) ContactsContext {
	return ContactsContext{Writer: w, Request: req, CorpID: corpID, Token: token, AesKey: aeskey}
}

// ParsePost 解析微信消息
func (c *ContactsContext) ParsePost() error {
	req := c.Request

	var eXML message.EncryptedXMLMsg
	err := xml.NewDecoder(req.Body).Decode(&eXML)
	if err != nil {
		return err
	}

	sign := req.FormValue("msg_signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")

	devSign := util.Signature(c.Token, timestamp, nonce, eXML.EncryptedMsg)
	if devSign != sign {
		return errors.New("sign failed")
	}

	_, rawMsg, err := util.DecryptMsg(c.CorpID, eXML.EncryptedMsg, c.AesKey)
	if err != nil {
		return err
	}

	c.Nonce = nonce
	log.App.Debug("contact:", string(rawMsg))
	return xml.Unmarshal(rawMsg, &c.Event)
}
