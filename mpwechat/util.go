package mpwechat

import (
	"net/http"

	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/util"
)

// AuthWxURL 微信消息接口验证
// 如果成功，返回 解密后的 字符串;不成功，返回空字符串
func AuthWxURL(req *http.Request, token string) string {
	log.App.Debug("auth get", req.URL.String())
	// msg_signature: 企业微信加密签名，msg_signature结合了企业填写的token、请求中的timestamp、nonce参数、加密的消息体
	// timestamp: 时间戳
	// nonce: 随机数
	// echostr: 加密的随机字符串
	// 1）将token、timestamp、nonce三个参数进行字典序排序 2）将三个参数字符串拼接成一个字符串进行sha1加密 3）开发者获得加密后的字符串可与signature对比，标识该请求来源于微信
	sign := req.FormValue("signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")
	echostr := req.FormValue("echostr")

	// 签名：dev_msg_signature=sha1(sort(token、timestamp、nonce、msg_encrypt))。sort的含义是将参数值按照字母字典排序，然后从小到大拼接成一个字符串
	devSign := util.Signature(token, timestamp, nonce)
	if devSign != sign {
		log.App.Debug("AuthWxURL failed:", devSign, sign)
		return ""
	}

	return echostr
}
