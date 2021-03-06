package corp

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/util"

	wxutil "github.com/kere/qywx/util"
)

const (
	ticketURL = "https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=%s"
)

// wxjs 验证的ticket
type ticketCached struct {
	cache.Map
}

func newTicketCached() *ticketCached {
	t := &ticketCached{}
	t.Init(t, 0)
	return t
}

// CheckValue 检查缓存值是否正确，如果正确缓存
func (t *ticketCached) CheckValue(v interface{}) bool {
	return v.(string) != ""
}

// Build func
func (t *ticketCached) Build(args ...interface{}) (interface{}, error) {
	corpName := args[0].(string)
	agentID := args[1].(int)

	cp, err := GetByName(corpName)
	if err != nil {
		return "", errors.New("corp not found in tokenCached")
	}
	agent := cp.GetAgentByAgentid(agentID)
	if agent == nil {
		return "", errors.New("agent not found in tokenCached")
	}

	token, err := agent.GetToken()
	if err != nil {
		return "", err
	}

	// 获取 access_ticket
	// 请求方式：GET（HTTPS）
	dat, err := wxutil.AjaxGet(fmt.Sprintf(ticketURL, token), nil)
	if err != nil {
		return "", err
	}
	t.SetExpires(dat.Int("expires_in"))

	return dat.String("ticket"), nil
}

// map[corpid] *ticketCached
var ticketCache = newTicketCached()

// GetTicket get cached token
func (a *Agent) GetTicket() (string, error) {
	ticket := ticketCache.Get(a.Corp.Name, a.Agentid)
	if ticket == nil {
		return "", errors.New("get cached ticket is nil")
	}
	return ticket.(string), nil
}

// WxAPIConfig p
func (a *Agent) WxAPIConfig(urlstr string) (util.MapData, error) {
	cp := a.Corp

	ticket, err := a.GetTicket()
	if err != nil {
		return nil, err
	}

	// signature
	nonce := wxutil.RandomStr(16)
	timestamp := fmt.Sprint(time.Now().Unix())

	arr := []string{
		"jsapi_ticket=" + ticket,
		"noncestr=" + nonce,
		"timestamp=" + timestamp,
		"url=" + urlstr,
	}

	h := sha1.New()
	io.WriteString(h, strings.Join(arr, "&"))
	sign := fmt.Sprintf("%x", h.Sum(nil))

	data := util.MapData{
		"signature": sign,
		"nonce":     nonce,
		"corpid":    cp.Corpid,
		"timestamp": timestamp,
		"url":       urlstr}

	return data, nil
}
