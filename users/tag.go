package users

import (
	"fmt"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
)

const (
	tagListURL = "https://qyapi.weixin.qq.com/cgi-bin/tag/list?access_token=%s"
	tagGetURL  = "https://qyapi.weixin.qq.com/cgi-bin/tag/get?access_token=%s&tagid=%d"
)

// Tag class
type Tag struct {
	ID   int    `json:"tagid"`
	Name string `json:"tagname"`
}

// GetTags tab list
func GetTags(token string) ([]Tag, error) {
	tags := make([]Tag, 0)

	dat, err := client.Get(fmt.Sprintf(tagListURL, token))
	if err != nil {
		return tags, err
	}

	log.App.Debug("get departments:", dat)

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

// GetTagUsers tab list
func GetTagUsers(tagid int, token string) (tagname string, usrs []User, pids []int, err error) {
	dat, err := client.Get(fmt.Sprintf(tagGetURL, token, tagid))
	if err != nil {
		return tagname, usrs, pids, err
	}

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
			pids = append(pids, int(d.(float64)))
		}
	}

	return dat.String("tagname"), usrs, pids, nil
}
