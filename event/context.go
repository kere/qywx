package event

// // ContactContext for message
// type ContactContext struct {
// 	Writer                http.ResponseWriter
// 	Request               *http.Request
// 	Nonce                 string
// 	CorpID, Token, AesKey string
//
// 	EventMessage ContactEvent
// }
//
// // NewContactContext func
// func NewContactContext(w http.ResponseWriter, req *http.Request, corpID, token, aeskey string) ContactContext {
// 	return ContactContext{Writer: w, Request: req, CorpID: corpID, Token: token, AesKey: aeskey}
// }
//
// // ParsePost 解析微信消息
// func (c *ContactContext) ParsePost() ([]byte, error) {
// 	req := c.Request
//
// 	var eXML message.EncryptedXMLMsg
// 	err := xml.NewDecoder(req.Body).Decode(&eXML)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	sign := req.FormValue("msg_signature")
// 	timestamp := req.FormValue("timestamp")
// 	nonce := req.FormValue("nonce")
//
// 	devSign := util.Signature(c.Token, timestamp, nonce, eXML.EncryptedMsg)
// 	if devSign != sign {
// 		return nil, errors.New("sign failed")
// 	}
//
// 	_, rawMsg, err := util.DecryptMsg(c.CorpID, eXML.EncryptedMsg, c.AesKey)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	c.Nonce = nonce
// 	log.App.Debug("contact:", string(rawMsg))
// 	return rawMsg, xml.Unmarshal(rawMsg, &c.EventMessage)
// }
