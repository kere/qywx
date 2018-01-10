package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

// Get unmarshal body and reture mapdata
func Get(uri string) (util.MapData, error) {
	return send("GET", uri, nil)
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

	log.App.Debug(uri)

	switch method {
	case "GET":
		body, err = get(uri)
	case "POST":
		body, err = post(uri, dat)
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

	if v.IsSet("errcode") && v.Int("errcode") > 0 {
		return nil, errors.New(v.String("errmsg"))
	}

	log.App.Debug(v)

	return v, nil
}

func get(uri string) ([]byte, error) {
	resq, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resq.Body.Close()
	return ioutil.ReadAll(resq.Body)
}

func post(uri string, dat util.MapData) ([]byte, error) {
	vals := url.Values{}
	for k, v := range dat {
		vals[k] = []string{fmt.Sprint(v)}
	}

	resq, err := http.PostForm(uri, vals)
	if err != nil {
		return nil, err
	}

	defer resq.Body.Close()
	return ioutil.ReadAll(resq.Body)
}

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
