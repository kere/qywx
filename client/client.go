package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

// Get unmarshal body and reture mapdata
func Get(uri string, dat util.MapData) (util.MapData, error) {
	return send("GET", uri, dat)
}

// PostForm unmarshal body and reture mapdata
func PostForm(uri string, dat util.MapData) (util.MapData, error) {
	return send("POST", uri, dat)
}

// PostJSON func
func PostJSON(uri string, dat util.MapData) (util.MapData, error) {
	return send("PostJson", uri, dat)
}

// send unmarshal body and reture mapdata
func send(method, uri string, dat util.MapData) (util.MapData, error) {
	var body []byte
	var err error

	log.App.Debug(method, uri)

	switch method {
	case "GET":
		body, err = util.AjaxGet(uri, dat)
	case "POST":
		body, err = util.AjaxPost(uri, dat)
	case "PostJson":
		body, err = postJSON(uri, dat)
	}
	if err != nil {
		return nil, err
	}

	var v util.MapData
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if v.IsSet("errcode") && v.Int("errcode") != 0 {
		return nil, errors.New(v.String("errmsg"))
	}

	log.App.Debug(method, v)

	return v, nil
}

// func get(uri string, dat util.MapData) ([]byte, error) {
// 	params := ""
// 	if dat != nil {
// 		vals := url.Values{}
// 		for k, v := range dat {
// 			vals.Add(k, fmt.Sprint(v))
// 		}
// 		params = "?" + vals.Encode()
// 	}
// 	resq, err := http.Get(uri + params)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	defer resq.Body.Close()
// 	return ioutil.ReadAll(resq.Body)
// }
//
// func post(uri string, dat util.MapData) ([]byte, error) {
// 	vals := url.Values{}
// 	for k, v := range dat {
// 		vals[k] = []string{fmt.Sprint(v)}
// 	}
//
// 	resq, err := http.PostForm(uri, vals)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	defer resq.Body.Close()
// 	return ioutil.ReadAll(resq.Body)
// }

func postJSON(uri string, dat util.MapData) ([]byte, error) {
	src, err := json.Marshal(dat)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
