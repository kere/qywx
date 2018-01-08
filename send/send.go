package send

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	utilib "github.com/kere/gno/libs/util"
	"github.com/kere/qywx/corp"
)

const (
	sendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

// Message class
type Message struct {
	Agent   *corp.Agent
	ToUser  []string
	ToParty []string
	ToTag   []string
	MsgType string
	Value   map[string]string
	Safe    int
}

// NewMessage func
func NewMessage(agent *corp.Agent, touser, topart, totab []string, msgtype string, val map[string]string) *Message {
	return &Message{Agent: agent, ToUser: touser, ToParty: topart, ToTag: totab, MsgType: msgtype, Value: val}
}

// Send text
func (m *Message) Send() (utilib.MapData, error) {
	touser := strings.Join(m.ToUser, "|")
	topart := strings.Join(m.ToParty, "|")
	totag := strings.Join(m.ToTag, "|")

	dat := utilib.MapData{
		"touser":  touser,
		"topart":  topart,
		"totag":   totag,
		"msgtype": m.MsgType,
		"agentid": m.Agent.ID}

	dat[m.MsgType] = m.Value

	token, err := m.Agent.GetToken()
	if err != nil {
		return nil, err
	}

	return Post(fmt.Sprintf(sendURL, token.Value), dat)
}

// Post func
func Post(uri string, dat utilib.MapData) (utilib.MapData, error) {
	src, err := json.Marshal(dat)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	src, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var obj utilib.MapData
	err = json.Unmarshal(src, obj)
	if err != nil {
		return nil, err
	}

	if obj.Int("errcode") > 0 {
		return obj, errors.New(obj.String("errmsg"))
	}
	return obj, nil
}
