package tag

import (
	"fmt"

	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
	"github.com/kere/qywx/depart"
)

const (
	tagListURL   = "https://qyapi.weixin.qq.com/cgi-bin/tag/list?access_token=%s"
	tagGetURL    = "https://qyapi.weixin.qq.com/cgi-bin/tag/get?access_token=%s&tagid=%d"
	tagCreateURL = "https://qyapi.weixin.qq.com/cgi-bin/tag/create?access_token=%s"
)

// Tag class
type Tag struct {
	ID   int    `json:"tagid"`
	Name string `json:"tagname"`
}

// WxAdd Create 标签, return tagid
func WxAdd(tagname, token string) (int, error) {
	dat, err := client.PostJSON(fmt.Sprint(tagCreateURL, token), util.MapData{"tagname": tagname})

	if err != nil {
		return 0, err
	}

	return dat.Int("tagid"), nil
}

// WxTags tab list
func WxTags(token string) ([]Tag, error) {
	tags := make([]Tag, 0)

	dat, err := client.Get(fmt.Sprintf(tagListURL, token), nil)
	if err != nil {
		return tags, err
	}

	// log.App.Debug("wxtags:", dat)

	arr := dat["taglist"].([]interface{})
	for _, d := range arr {
		item := util.MapData(d.(map[string]interface{}))
		tags = append(tags, Tag{
			ID:   item.Int("tagid"),
			Name: item.String("tagname"),
		})
	}

	return tags, nil
}

// User class
type User struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
}

// WxTagUsers 只列出标签，不抓取部门下的用户
func WxTagUsers(tagid int, token string) (tagname string, usrs []User, partyIds []int, err error) {
	dat, err := client.Get(fmt.Sprintf(tagGetURL, token, tagid), nil)
	if err != nil {
		return tagname, usrs, partyIds, err
	}
	// log.App.Debug("wxtag users:", dat)

	if arr, isok := dat["userlist"].([]interface{}); isok {
		for _, d := range arr {
			item := util.MapData(d.(map[string]interface{}))
			usrs = append(usrs, User{
				UserID: item.String("userid"),
				Name:   item.String("name"),
			})
		}
	}

	if arr, isok := dat["partylist"].([]interface{}); isok {
		for _, d := range arr {
			partyIds = append(partyIds, int(d.(float64)))
		}
	}

	return dat.String("tagname"), usrs, partyIds, nil
}

// WxTagFullUsers 包括部门下的所有用户信息
func WxTagFullUsers(tagid int, token string) (tagname string, usrs []User, partyIds []int, dusrs []depart.User, useridlist []string, err error) {
	tagname, usrs, partyIds, err = WxTagUsers(tagid, token)
	if err != nil {
		return tagname, usrs, partyIds, dusrs, useridlist, err
	}

	var arr []depart.User
	for i := range partyIds {
		arr, err = depart.WxDepartSimpleUsers(partyIds[i], true, token)
		if err != nil {
			return tagname, usrs, partyIds, dusrs, useridlist, err
		}
		dusrs = append(dusrs, arr...)
	}

	// useridlist = []string{}
	for _, u := range usrs {
		useridlist = append(useridlist, u.UserID)
	}
	for _, u := range dusrs {
		useridlist = append(useridlist, u.UserID)
	}

	return tagname, usrs, partyIds, dusrs, useridlist, err
}
