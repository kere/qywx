package users

import (
	"fmt"

	"github.com/kere/qywx/client"
)

const (
	userGetURL = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
)

// UserDetail class
type UserDetail struct {
	ID          string `json:"userid"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Department  []int  `json:"department"` // 部门
	Position    string `json:"position"`   // 职位
	// Order  []int  `json:"order"` // 部门排序
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"` // 性别
	Email    string `json:"email"`
	IsLeader int    `json:"isleader"`
	Avatar   string `json:"avatar"`
	Tele     string `json:"telephone"`
	Status   int    `json:"status"`
	Enable   int    `json:"enable"`
	// ExtAttr  map[string]string `json:"extattr"`
}

// WxUser func
func WxUser(userID, token string) (UserDetail, error) {
	var usr UserDetail
	dat, err := client.Get(fmt.Sprintf(userGetURL, token, userID))
	if err != nil {
		return usr, err
	}

	usr.ID = dat.String("userid")
	usr.Name = dat.String("name")
	usr.EnglishName = dat.StringDefault("english_name", "")
	usr.Department = dat.IntsDefault("department", []int{})
	// usr.Order = dat.String("order")
	usr.Position = dat.StringDefault("position", "")
	usr.Mobile = dat.StringDefault("mobile", "")
	usr.Email = dat.StringDefault("email", "")
	usr.Avatar = dat.StringDefault("avatar", "")
	usr.Tele = dat.StringDefault("telephone", "")
	usr.IsLeader = dat.IntDefault("isleader", 0)
	usr.Gender = dat.String("gender")
	usr.Status = dat.IntDefault("status", 0)
	usr.Enable = dat.IntDefault("enable", 0)

	return usr, nil
}
