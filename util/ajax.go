package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

// AjaxGet unmarshal body and reture mapdata
func AjaxGet(uri string, vals url.Values) (util.MapData, error) {
	body, err := util.AjaxGet(uri, vals)
	if err != nil {
		return nil, err
	}
	log.App.Debug(uri)
	return parseBody(body)
}

// PostForm unmarshal body and reture mapdata
func PostForm(uri string, dat util.MapData) (util.MapData, error) {
	body, err := util.AjaxPost(uri, dat)
	if err != nil {
		return nil, err
	}
	log.App.Debug(uri)
	return parseBody(body)
}

// PostJSON func
func PostJSON(uri string, dat interface{}) (util.MapData, error) {
	body, err := postJSON(uri, dat)
	if err != nil {
		return nil, err
	}
	log.App.Debug(uri)
	return parseBody(body)
}

// parseBody unmarshal body and reture mapdata
func parseBody(body []byte) (util.MapData, error) {
	var v util.MapData
	err := json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	if v.IsSet("errcode") && v.Int("errcode") != 0 {
		return nil, errors.New(v.String("errcode") + ":" + v.String("errmsg"))
	}

	return v, nil
}

func postJSON(uri string, dat interface{}) ([]byte, error) {
	src, err := json.Marshal(dat)
	if err != nil {
		return nil, err
	}

	src = bytes.Replace(src, []byte("\\u003c"), []byte("<"), -1)
	src = bytes.Replace(src, []byte("\\u003e"), []byte(">"), -1)
	src = bytes.Replace(src, []byte("\\u0026"), []byte("&"), -1)

	resp, err := http.Post(uri, "application/json;charset=utf-8", bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
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
