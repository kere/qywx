package mpwechat

var (
	wxAPIarea = ""
	apihttp   = "https://"
)
var apiurls = map[string]string{
	"Token":      "api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
	"MenuCreate": "api.weixin.qq.com/cgi-bin/menu/create?access_token=%s",
	"MenuGet":    "api.weixin.qq.com/cgi-bin/menu/get?access_token=%s",
	"MenuDelete": "api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s",
}

// SetAPIArea 设置API的区域
func SetAPIArea(area string) {
	wxAPIarea = area + "."
}

// WxURL 微信接口API
func WxURL(name string) string {
	v, isok := apiurls[name]
	if isok {
		return apihttp + wxAPIarea + v
	}
	panic("WxURL not found:" + name)
}
