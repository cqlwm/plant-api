package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var app_id = "wx579832410b69b5a0"
var secret = "2341eed48ccd59c5f6520d927dc9bcb6"

// auth.code2Session
var tagWeiXinAuth = "WeiXinAuth ## "

func WeiXinAuth(jsCode string) (id string, key string, b bool) {
	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		app_id, secret, jsCode)

	resp, err := http.Get(wxUrl)
	if err != nil {
		log.Println(tagWeiXinAuth, err)
		return "", "", false
	}
	if resp.StatusCode != 200 {
		log.Println(tagWeiXinAuth, "请求失败")
		return "", "", false
	}

	result := wxNetResult{}
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(tagWeiXinAuth, err)
		return "", "", false
	}

	err = json.Unmarshal(rb, &result)
	if err != nil {
		log.Println(tagWeiXinAuth, err)
		return "", "", false
	}
	if result.Errcode != 0 {
		return string(result.Errcode), result.Errmsg, false
	}

	id = result.Openid
	key = result.Session_key
	return id, key, true
}

/*
	openid	string	用户唯一标识
	session_key	string	会话密钥
	unionid	string	用户在开放平台的唯一标识符，若当前小程序已绑定到微信开放平台帐号下会返回，详见 UnionID 机制说明。
	errcode	number	错误码
	errmsg	string	错误信息
*/
type wxNetResult struct {
	Openid      string
	Session_key string
	Unionid     string
	Errcode     int
	Errmsg      string
}
