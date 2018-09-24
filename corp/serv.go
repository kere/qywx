package corp

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/message"
	"github.com/kere/qywx/util"
)

// IExec interface
type IExec interface {
	IsExec(message.Context) bool
	Exec(message.Context, httprouter.Params) message.IMessage
}

// Serv message
type Serv struct {
	Execs []IExec
}

// NewServ func
func NewServ() *Serv {
	return &Serv{}
}

//AddExec 添加处理模块
func (srv *Serv) AddExec(exec IExec) {
	srv.Execs = append(srv.Execs, exec)
}

//Auth 验证
func (srv *Serv) Auth(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	cp, agent, err := getCorpAndAgent(ps)
	if err != nil {
		return
	}
	s := util.AuthWxURL(req, cp.Corpid, agent.MsgToken, agent.MsgAesKey)
	rw.Write([]byte(s))
}

//MessageHandle 处理微信的请求
func (srv *Serv) MessageHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	cp, agent, err := getCorpAndAgent(ps)
	if err != nil {
		return
	}

	ctx := message.NewContext(rw, req, cp.Corpid, agent.Secret, agent.MsgToken, agent.MsgAesKey)
	ctx.IsSafe = true
	srv.doExec(ctx, ps)
}

func (srv *Serv) doExec(ctx message.Context, ps httprouter.Params) error {
	_, err := ctx.ParsePost()
	if err != nil {
		return err
	}
	log.App.Debug("message:", ctx.MixMessage.GetMsgType(), ctx.MixMessage.Event, ctx.MixMessage.Content)
	// exec replies
	for _, exec := range srv.Execs {
		if exec.IsExec(ctx) {
			msg := exec.Exec(ctx, ps)
			return ctx.Send(msg)
		}
	}
	return nil
}

func getCorpAndAgent(ps httprouter.Params) (*Corporation, *Agent, error) {
	cp, err := GetByName(ps.ByName("corp"))
	if err != nil {
		return nil, nil, err
	}
	agentName := ps.ByName("agent")
	if agentName == "" {
		return cp, nil, nil
	}
	agent, err := cp.GetAgent(agentName)
	if err != nil {
		return cp, nil, err
	}
	return cp, agent, nil
}

// ContactsServ message
type ContactsServ struct {
	Serv
}

// NewContaServ func
func NewContaServ() *ContactsServ {
	return &ContactsServ{}
}

//Auth 验证
func (srv *ContactsServ) Auth(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	cp, _, err := getCorpAndAgent(ps)
	if err != nil {
		return
	}
	s := util.AuthWxURL(req, cp.Corpid, cp.ContactsToken, cp.ContactsAesKey)
	rw.Write([]byte(s))
}

//MessageHandle 处理微信的请求
func (srv *ContactsServ) MessageHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	cp, _, err := getCorpAndAgent(ps)
	if err != nil {
		return
	}

	ctx := message.NewContext(rw, req, cp.Corpid, cp.ContactsSecret, cp.ContactsToken, cp.ContactsAesKey)
	ctx.IsSafe = true
	srv.doExec(ctx, ps)
	log.App.Debug(ctx.MixMessage.GetMsgType(), ctx.MixMessage.Event, ctx.MixMessage.Content)
}
