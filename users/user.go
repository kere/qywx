package users

import (
	"fmt"

	"github.com/kere/gno/libs/log"
	"github.com/kere/qywx/client"
)

const (
	userGetURL = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
)

// UserDetail class
type UserDetail struct {
	ID          string            `json:"userid"`
	Name        string            `json:"name"`
	EnglishName string            `json:"english_name"`
	Department  []int             `json:"department"`
	Position    string            `json:"position"` // 职位
	Order       string            `json:"order"`
	Mobile      string            `json:"mobile"`
	Gender      int               `json:"gender"` // 性别
	Email       string            `json:"email"`
	IsLeader    string            `json:"isleader"`
	Avatar      string            `json:"avatar"`
	Tele        string            `json:"telephone"`
	Status      int               `json:"status"`
	ExtAttr     map[string]string `json:"extattr"`
}

// GetUser func
func GetUser(userID, token string) (UserDetail, error) {
	var usr UserDetail
	dat, err := client.GetMapData(fmt.Sprintf(userGetURL, token, userID))
	if err != nil {
		return usr, err
	}
	log.App.Debug("get wx user:", dat)

	usr.ID = dat.String("userid")
	usr.Name = dat.String("name")
	usr.EnglishName = dat.String("english_name")
	usr.Department = dat.Ints("depatment")
	usr.Order = dat.String("order")
	usr.Position = dat.String("position")
	usr.Mobile = dat.String("mobile")
	usr.Email = dat.String("email")
	usr.Avatar = dat.String("avatar")
	usr.Tele = dat.String("telephone")
	usr.Gender = dat.Int("gender")
	usr.Status = dat.Int("status")

	return usr, nil
}
