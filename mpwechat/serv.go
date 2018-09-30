package mpwechat

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/message"
)

// IExec interface
type IExec interface {
	IsExec(message.Context) bool
	Exec(message.Context, httprouter.Params) message.IMessage
}

// Serv message
type Serv struct {
	AppID     string
	AppSecret string
	Token     string
	AesKey    string

	IsSafe bool

	Execs []IExec
}

// NewServ func
func NewServ(appid, appsecret, token, aeskey string) *Serv {
	return &Serv{AppID: appid, AppSecret: appsecret, Token: token, AesKey: aeskey}
}

//AddExec 添加处理模块
func (srv *Serv) AddExec(exec IExec) {
	srv.Execs = append(srv.Execs, exec)
}

//Auth 验证
func (srv *Serv) Auth(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
	for _, exec := range srv.Execs {
		if exec.IsExec(ctx) {
			msg := exec.Exec(ctx, ps)
			if msg.GetToUserName() == "" {
				break
			}
			if err := ctx.Send(msg); err != nil {
				log.App.Error(err)
			}
		}
	}

}
