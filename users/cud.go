package users

import (
	"errors"
	"fmt"

	"github.com/kere/gno/libs/util"
	wxutil "github.com/kere/qywx/util"
)

const (
	userCreateURL = "https://qyapi.weixin.qq.com/cgi-bin/user/create?access_token=%s"
	userUpdateURL = "https://qyapi.weixin.qq.com/cgi-bin/user/update?access_token=%s"
	userDeleteURL = "https://qyapi.weixin.qq.com/cgi-bin/user/delete?access_token=%s&userid=%s"
)

// WxCreate user
func WxCreate(dat util.MapData, token string) error {
	if dat.String("userid") == "" {
		return errors.New("UserID不能为空")
	}
	if dat.String("name") == "" {
		return errors.New("用户名称不能为空")
	}
	if dat.String("mobile") == "" {
		return errors.New("用户手机不能为空")
	}
	if len(dat.Ints("department")) == 0 {
		return errors.New("部门不能为空")
	}

	_, err := wxutil.PostJSON(fmt.Sprintf(userCreateURL, token), dat)
	return err
}

// WxUpdate user
func WxUpdate(dat util.MapData, token string) error {
	if dat.String("userid") == "" {
		return errors.New("修改用户的UserID不能为空")
	}
	if dat.String("name") == "" {
		return errors.New("修改用户的用户名称不能为空")
	}
	if dat.String("mobile") == "" {
		return errors.New("修改用户的用户手机不能为空")
	}
	if len(dat.Ints("department")) == 0 {
		return errors.New("修改用户的部门不能为空")
	}

	_, err := wxutil.PostJSON(fmt.Sprintf(userUpdateURL, token), dat)
	return err
}

// WxDelete user
func WxDelete(userid, token string) error {
	if userid == "" {
		return errors.New("删除操作的用户的UserID不能为空")
	}

	_, err := wxutil.AjaxGet(fmt.Sprintf(userDeleteURL, token, userid), nil)
	return err
}
