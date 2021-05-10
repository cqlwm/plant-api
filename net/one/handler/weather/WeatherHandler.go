package weather

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"plant-api/net/one/config"
	"plant-api/net/one/crawling"
	"plant-api/net/one/service/impl/weather"
	"strconv"
	"time"
)

// handler路由配置
func WeatherHandler(e *gin.Engine) {
	e.GET(config.WeatherHandler2query, query)
	e.GET(config.WeatherHandler2simple, simple)
}

var redisTTL time.Duration = 1800
var ws *weather.WeatherServiceImpl

func query(c *gin.Context) {
	code := c.Query("code")
	fmt.Println(code)

	i, err := strconv.Atoi(code)
	if err != nil {
		config.Error(c, &config.PARAM_ERROR, nil)
		return
	}

	infoOne := crawling.WeatherResult{}

	err = config.GetJson(config.WeatherRedisKey+code, &infoOne)
	if err == nil && infoOne.Info != nil {
		fmt.Println("我走了缓存")
		config.Ok2(c, infoOne)
		return
	}

	info, err := ws.Info(i)
	if err != nil {
		config.Error2(c, err.Error())
		return
	}

	_ = config.SaveKJson(config.WeatherRedisKey+code, info, redisTTL)

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

	infoOne := make(map[string]interface{}, 0)

	err = config.GetJson(config.SimpleWeatherRedisKey+code, &infoOne)
	if err == nil && len(infoOne) != 0 {
		fmt.Println("我走了缓存")
		config.Ok2(c, infoOne)
		return
	}

	info := ws.Simple(iCode)

	_ = config.SaveKJson(config.SimpleWeatherRedisKey+code, info, redisTTL)

	config.Ok2(c, info)
}
