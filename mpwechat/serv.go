package mpwechat

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/message"
)

// IReply interface
type IReply interface {
	Match(msg *message.Context, txt string) (bool, [][]string)
	Build(*message.Context, [][]string) (message.IMessage, error)
}

// Serv message
type Serv struct {
	AppID     string
	AppSecret string
	Token     string
	AesKey    string

	IsSafe bool

	Replies []IReply
}

// NewServ func
func NewServ(appid, appsecret, token, aeskey string) *Serv {
	return &Serv{AppID: appid, AppSecret: appsecret, Token: token, AesKey: aeskey}
}

//AddReply 添加处理模块
func (srv *Serv) AddReply(exec IReply) {
	srv.Replies = append(srv.Replies, exec)
}

//AuthWx 验证
func (srv *Serv) AuthWx(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s := AuthWxURL(req, srv.Token)
	rw.Write([]byte(s))
}

//MessageHandle 处理微信的请求
func (srv *Serv) MessageHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ctx := message.NewContext(rw, req, srv.AppID, srv.AppSecret, srv.Token, srv.AesKey)
	ctx.IsSafe = srv.IsSafe

	_, err := ctx.ParsePost()
	if err != nil {
		log.App.Error(err)
		return
	}

	log.App.Debug(ctx.MixMessage.GetMsgType(), ctx.MixMessage.Event, ctx.MixMessage.EventKey, ctx.MixMessage.Content)

	// exec replies
	for _, exec := range srv.Replies {
		isDo, msg, err := RunReply(&ctx, exec)
		if !isDo {
			continue
		}
		if err != nil {
			log.App.Error(err)
			ctx.Send(message.NewReplyText("写代码的程序员犯了一个错误，已经后台通知管理员解决", &ctx.MixMessage))
			break
		}

		if err := ctx.Send(msg); err != nil {
			log.App.Error(err)
		}
	}

}

func trimStr(r rune) bool {
	return r == ' ' || r == '\n'
}

// RunReply reply
func RunReply(ctx *message.Context, reply IReply) (bool, message.IMessage, error) {
	msg := &ctx.MixMessage
	content := strings.TrimFunc(msg.Content, trimStr)
	isM, result := reply.Match(ctx, content)
	if !isM {
		return false, message.NewEmpty(), nil
	}

	m, err := reply.Build(ctx, result)
	return true, m, err
}
