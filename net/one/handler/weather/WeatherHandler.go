package weather

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"plant-api/net/one/config"
	"plant-api/net/one/service/impl/weather"
	"strconv"
)

// handler路由配置
func WeatherHandler(e *gin.Engine) {
	e.GET(config.WeatherHandler2query, query)
	e.GET(config.WeatherHandler2simple, simple)
}

var ws *weather.WeatherServiceImpl

func query(c *gin.Context) {
	code := c.Query("code")
	fmt.Println(code)

	i, err := strconv.Atoi(code)
	if err != nil {
		config.Error(c, &config.PARAM_ERROR, nil)
		return
	}

	info, err := ws.Info(i)
	if err != nil {
		config.Error2(c, err.Error())
		return
	}
	config.Ok2(c, info)
}

func simple(c *gin.Context) {
	code := c.Query("code")
	fmt.Println(code)

	iCode, err := strconv.Atoi(code)
	if err != nil {
		config.Error(c, &config.PARAM_ERROR, nil)
		return
	}

	info := ws.Simple(iCode)
	if err != nil {
		config.Error2(c, err.Error())
		return
	}
	config.Ok2(c, info)
}
