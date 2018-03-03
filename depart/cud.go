package depart

import (
	"errors"
	"fmt"

	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
)

const (
	departCreateURL = "https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token=%s"
	departUpdateURL = "https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token=%s"
	departDelteURL  = "https://qyapi.weixin.qq.com/cgi-bin/department/delete?access_token=%s&id=%d"
)

// WxCreate department
func WxCreate(name string, parentID, order int, token string) (int, error) {
	if name == "" {
		return -1, errors.New("创建部门的名称不能为空")
	}

	args := util.MapData{"name": name, "parent_id": parentID, "order": order}
	dat, err := client.PostJSON(fmt.Sprintf(departCreateURL, token), args)

	if err != nil {
		return -1, err
	}
	// log.App.Debug("wxdeparts:", dat)

	return dat.Int("id"), nil
}

// WxUpdate department
func WxUpdate(id int, name string, parentID, order int, token string) error {
	if name == "" {
		return errors.New("更新部门的名称不能为空")
	}

	args := util.MapData{"id": id, "name": name, "parent_id": parentID, "order": order}
	_, err := client.PostJSON(fmt.Sprintf(departUpdateURL, token), args)
	return err
}

// WxDelete department
func WxDelete(id int, token string) error {
	if id == 0 {
		return errors.New("删除部门ID等于0")
	}

	_, err := client.Get(fmt.Sprintf(departDelteURL, token, id))
	return err
}
