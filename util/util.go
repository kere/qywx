package util

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/kere/gno"
	"github.com/kere/gno/libs/log"
)

var (
	xmlContentType   = []string{"text/xml; charset=utf-8"}
	plainContentType = []string{"text/plain; charset=utf-8"}
	expiresVal       = 0
)

// Expires int
func Expires() int {
	if expiresVal > 0 {
		return expiresVal
	}

	expiresVal = gno.GetConfig().GetConf("data").DefaultInt("data_expires", 72000)
	return expiresVal
}

//Signature sha1签名
func Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// WriteContextType func
func WriteContextType(w http.ResponseWriter, headers []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = headers
	}
}

// AuthWxURL 微信消息接口验证
// 如果成功，返回 解密后的 字符串;不成功，返回空字符串
func AuthWxURL(req *http.Request, corpID, token, aeskey string) string {
	log.App.Debug("auth get", req.URL.String())
	// msg_signature: 企业微信加密签名，msg_signature结合了企业填写的token、请求中的timestamp、nonce参数、加密的消息体
	// timestamp: 时间戳
	// nonce: 随机数
	// echostr: 加密的随机字符串，以msg_encrypt格式提供。需要解密并返回echostr明文，解密后有random、msg_len、msg、$CorpID四个字段，其中msg即为echostr明文
	sign := req.FormValue("msg_signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")
	echostr := req.FormValue("echostr")
	// 企业计算签名：dev_msg_signature=sha1(sort(token、timestamp、nonce、msg_encrypt))。sort的含义是将参数值按照字母字典排序，然后从小到大拼接成一个字符串
	devSign := Signature(token, timestamp, nonce, echostr)
	if devSign != sign {
		log.App.Debug("AuthWxURL:", devSign, sign)
		return ""
	}

	_, rawMsg, err := DecryptMsg(corpID, echostr, aeskey)
	if err != nil {
		log.App.Debug("AuthWxURL:", err)
		return ""
	}

	return fmt.Sprintf("%s", rawMsg)
}

// WriteXML 写入xml信息
func WriteXML(w http.ResponseWriter, obj interface{}) {
	if obj == nil {
		return
	}

	WriteContextType(w, xmlContentType)
	bytes, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	_, err = w.Write(bytes)
	if err != nil {
		panic(err)
	}
}

// WriteString 写入string信息
func WriteString(w http.ResponseWriter, txt string) {
	WriteContextType(w, plainContentType)
	w.WriteHeader(200)
	_, err := w.Write([]byte(txt))
	if err != nil {
		panic(err)
	}
}
