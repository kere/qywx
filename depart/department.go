package depart

import (
	"fmt"

	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
)

const (
	departGetURL   = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=%s"
	departUsersURL = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=%s&department_id=%d&fetch_child=%d"
)

// Department class
type Department struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parentid"`
	Order    int    `json:"order"`
}

// WxDepartments func
func WxDepartments(departID int, token string) ([]Department, error) {
	dptList := make([]Department, 0)
	s := ""
	if departID > 0 {
		s = fmt.Sprint(departID)
	}

	dat, err := client.Get(fmt.Sprintf(departGetURL, token, s))
	if err != nil {
		return dptList, err
	}
	// log.App.Debug("wxdeparts:", dat)

	arr := dat["department"].([]interface{})
	for _, d := range arr {
		dpt := util.MapData(d.(map[string]interface{}))
		dptList = append(dptList, Department{
			ID:       dpt.Int("id"),
			Name:     dpt.String("name"),
			ParentID: dpt.Int("parentid"),
			Order:    dpt.Int("order"),
		})
	}

	return dptList, nil
}

// User class
type User struct {
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
}

// WxDepartUsers func
func WxDepartUsers(departID int, isChild bool, token string) ([]User, error) {
	usrs := make([]User, 0)
	isChildV := 0
	if isChild {
		isChildV = 1
	}

	dat, err := client.Get(fmt.Sprintf(departUsersURL, token, departID, isChildV))
	if err != nil {
		return usrs, err
	}
	// log.App.Debug("wxdepart users:", dat)

	arr := dat["userlist"].([]interface{})
	for _, d := range arr {
		dpt := util.MapData(d.(map[string]interface{}))
		usrs = append(usrs, User{
			UserID:     dpt.String("userid"),
			Name:       dpt.String("name"),
			Department: dpt.Ints("department"),
		})
	}

	return usrs, nil
}
