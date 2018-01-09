package send

import (
	"fmt"
	"strings"

	utilib "github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
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

	return client.PostJSON(fmt.Sprintf(sendURL, token.Value), dat)
}
