package user

import (
	"github.com/gin-gonic/gin"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
	"plant-api/net/one/service/impl/user"
)

// handler路由配置
func UserHandler(e *gin.Engine) {
	e.POST(config.UserHandler2Login, login)
}

var usi = user.UserServiceImpl{}

// 微信登录
func login(c *gin.Context) {
	form := entry.LoginForm{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		config.Error(c, &config.PARAM_ERROR, nil)
		return
	}

	// login check
	token, err := usi.LoginCheck(&form)
	if err != nil {
		config.Error2(c, err.Error())
		return
	}
	config.Ok2(c, token)
}
