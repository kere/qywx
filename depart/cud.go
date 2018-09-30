package depart

import (
	"errors"
	"fmt"

	"github.com/kere/gno/libs/util"
	wxutil "github.com/kere/qywx/util"
)

const (
	departCreateURL = "https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token=%s"
	departUpdateURL = "https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token=%s"
	departDelteURL  = "https://qyapi.weixin.qq.com/cgi-bin/department/delete?access_token=%s&id=%d"
)

// WxCreate department
func WxCreate(name string, parentID int, order int32, token string) (int, error) {
	if name == "" {
		return -1, errors.New("创建部门的名称不能为空")
	}

	args := util.MapData{"name": name, "parentid": parentID}
	if order > 0 {
		args["order"] = order
	}

	dat, err := wxutil.PostJSON(fmt.Sprintf(departCreateURL, token), args)

	if err != nil {
		return -1, err
	}
	// log.App.Debug("wxdeparts:", dat)
	ClearDepart()

	return dat.Int("id"), nil
}

// WxUpdateName department
func WxUpdateName(id int, name string, token string) error {
	if name == "" {
		return errors.New("更新部门的名称不能为空")
	}

	dat := util.MapData{"id": id, "name": name}
	return WxUpdate(id, dat, token)
}

// WxUpdate department
func WxUpdate(id int, dat util.MapData, token string) error {
	_, err := wxutil.PostJSON(fmt.Sprintf(departUpdateURL, token), dat)

	ClearDepart()
	return err
}

// WxDelete department
func WxDelete(id int, token string) error {
	if id == 0 {
		return errors.New("删除部门ID等于0")
	}

	_, err := wxutil.AjaxGet(fmt.Sprintf(departDelteURL, token, id), nil)
	ClearDepart()
	return err
}
