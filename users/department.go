package users

import (
	"fmt"

	"github.com/kere/gno/libs/util"
	"github.com/kere/qywx/client"
)

const (
	departGetURL = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=%s"
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
