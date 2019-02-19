package send

import (
	"fmt"
	"strings"

	utilib "github.com/kere/gno/libs/util"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/util"
)

const (
	sendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

// SendText text
// ToUser  []string 成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向该企业应用的全部成员发送
// ToParty []string
// ToTag   []string
func SendText(agent *corp.Agent, txt string, toUsers, toParty, toTags []string) (utilib.MapData, error) {
	touser := strings.Join(toUsers, "|")
	topart := strings.Join(toParty, "|")
	totag := strings.Join(toTags, "|")

	dat := utilib.MapData{
		"touser":  touser,
		"topart":  topart,
		"totag":   totag,
		"msgtype": "text",
		"text": utilib.MapData{
			"content": txt,
		},
		"agentid": agent.Agentid}

	token, err := agent.GetToken()
	if err != nil {
		return nil, err
	}

	return util.PostJSON(fmt.Sprintf(sendURL, token), dat)
}
