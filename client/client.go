package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

// GetMapData unmarshal body and reture mapdata
func GetMapData(uri string) (util.MapData, error) {
	return send("GET", uri, nil)
}

// PostMapData unmarshal body and reture mapdata
func PostMapData(uri string, dat util.MapData) (util.MapData, error) {
	return send("POST", uri, dat)
}

// send unmarshal body and reture mapdata
func send(method, uri string, dat util.MapData) (util.MapData, error) {
	var body []byte
	var err error
	switch method {
	case "GET":
		body, err = get(uri)
	case "POST":
		body, err = post(uri, dat)
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

	return v, nil
}

func get(uri string) ([]byte, error) {
	resq, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	log.App.Debug(uri)
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
	log.App.Debug(uri, dat)

	defer resq.Body.Close()
	return ioutil.ReadAll(resq.Body)
}
