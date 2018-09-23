package mpwechat

import (
	"errors"
	"fmt"

	"github.com/kere/gno/db"
	"github.com/kere/qywx/client"
)

// MenuGet create
func MenuGet(appid, appsecret string) (db.DataSet, error) {
	token, err := WxAccessToken(appid, appsecret)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(WxURL("MenuGet"), token)
	dat, err := client.Get(uri, nil)
	if err != nil {
		return nil, err
	}
	if !dat.IsSet("menu") {
		return nil, errors.New("no field menu found")
	}
	menu := dat["menu"].(map[string]interface{})
	btns := menu["button"].([]interface{})
	arr := db.DataSet{}
	for _, btn := range btns {
		arr = append(arr, btn.(map[string]interface{}))
	}

	fmt.Println(arr)
	return arr, nil
}

// MenuCreate create
func MenuCreate(appid, appsecret string, opt interface{}) error {
	token, err := WxAccessToken(appid, appsecret)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf(WxURL("MenuCreate"), token)
	_, err = client.PostJSON(uri, nil)
	return err
}

// MenuDelete create
func MenuDelete(appid, appsecret string) error {
	token, err := WxAccessToken(appid, appsecret)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf(WxURL("MenuDelete"), token)
	_, err = client.Get(uri, nil)
	return err
}
