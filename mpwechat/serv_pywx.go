package mpwechat

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/corp"
	"github.com/kere/qywx/message"
	"github.com/kere/qywx/util"
)

// IReplyQYWX interface
type IReplyQYWX interface {
	Regexp() *regexp.Regexp
	Build(*corp.Agent, [][]string) (message.IMessage, error)
}

// ServPYWX class
type ServPYWX struct {
	IsSafe bool

	Replies                       []IReplyQYWX
	Corpid, Secret, Token, AesKey string
	Agent                         *corp.Agent
}

// NewServPYWX new
func NewServPYWX(replys []IReplyQYWX) *ServPYWX {
	m := &ServPYWX{}
	m.IsSafe = true
	for i := range replys {
		m.AddReply(replys[i])
	}
	return m
}

//AddReply 添加处理模块
func (srv *ServPYWX) AddReply(exec IReplyQYWX) {
	srv.Replies = append(srv.Replies, exec)
}

// Prepare 验证微信消息地址
func (srv *ServPYWX) Prepare(ps httprouter.Params) bool {
	cp, err := corp.GetByName(ps.ByName("corp"))
	if err != nil {
		return false
	}
	agentName := ps.ByName("agent")
	if agentName == "" {
		return false
	}
	agent, err := cp.GetAgent(agentName)
	if err != nil {
		return false
	}
	srv.Agent = agent
	srv.Corpid = cp.Corpid
	srv.Secret = agent.Secret
	srv.Token = agent.MsgToken
	srv.AesKey = agent.MsgAesKey
	return true
}

// AuthWx 验证微信消息地址
func (srv *ServPYWX) AuthWx(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !srv.Prepare(ps) {
		return
	}
	str := util.AuthWxURL(req, srv.Corpid, srv.Token, srv.AesKey)
	rw.Write([]byte(str))
}

//MessageHandle 处理微信的请求
func (srv *ServPYWX) MessageHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if !srv.Prepare(ps) {
		return
	}
	// ctx := message.NewContext(rw, req, srv.AppID, srv.AppSecret, srv.Token, srv.AesKey)
	ctx := message.NewContext(rw, req, srv.Corpid, srv.Secret, srv.Token, srv.AesKey)
	ctx.IsSafe = srv.IsSafe

	_, err := ctx.ParsePost()
	if err != nil {
		log.App.Error(err)
		return
	}

	// exec replies
	isRun := false
	for _, exec := range srv.Replies {
		msg := ctx.MixMessage
		content := strings.TrimFunc(msg.Content, trimStr)

		result := exec.Regexp().FindAllStringSubmatch(content, -1)
		log.App.Debug("Message:", content, " Result:", result)

		if len(result) == 0 || len(result[0]) < 2 {
			continue
		}

		rMsg, err := exec.Build(srv.Agent, result)

		if err != nil {
			log.App.Error(err)
			ctx.Send(message.NewReplyText("写代码的程序员犯了一个错误，已经后台通知管理员解决", &ctx.MixMessage))
			break
		}

		if err := ctx.Send(rMsg); err != nil {
			log.App.Error(err)
		}
		isRun = true
	}
	if !isRun {
		ctx.Send(message.NewReplyText("是不是名称错了？没有找到任何答案！", &ctx.MixMessage))
	}
}
