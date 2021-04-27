package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"plant-api/net/one/config"
	"plant-api/net/one/crawling"
	"plant-api/net/one/entry"
	"plant-api/net/one/service/impl/article"
	"plant-api/net/one/service/impl/logs"
)

// handler路由配置
func TestHandlerRoute(e *gin.Engine) {
	e.GET(config.TestHandlerPing, ping)
	e.GET("/plant-api/ILogServiceTest", ILogServiceTest)
	e.GET("/plant-api/InterfaceTest", InterfaceTest)
	e.GET("/plant-api/LogCountTest", LogCountTest)
	e.GET("/plant-api/IntervalAskUrlTest", IntervalAskUrlTest)
	e.GET("/plant-api/wzcs", wzcs)
}

// ping测试程序是否启动
func ping(c *gin.Context) {
	config.Ok(c, &config.SUCCESS, "程序已经启动...")
}

// interface test handler
func InterfaceTest(c *gin.Context) {
	i := map[string]interface{}{}

	i["hello"] = "world"
	i["string"] = entry.IdentifyLog{
		Id:       0,
		UserId:   0,
		Option:   "333",
		OptionId: 0,
		Content:  "666",
		Created:  0,
	}

	var ret interface{} = i

	config.Ok2(c, ret)
}

// 日志测试
func ILogServiceTest(c *gin.Context) {
	log := entry.IdentifyLog{
		UserId:   100,
		Option:   "weeds",
		OptionId: 100,
	}

	impl := logs.IdentifyLogServiceImpl{}
	impl.Log(&log)
	config.Ok2(c, nil)
}

// 日志统计测试
func LogCountTest(c *gin.Context) {
	ils := logs.IdentifyLogServiceImpl{}
	count := ils.LogCount(100)
	config.Ok2(c, count)
}

// 文章爬虫测试
func IntervalAskUrlTest(c *gin.Context) {
	crawling.IntervalAskUrl()
}

// 文章测试
func wzcs(c *gin.Context) {
	i := []int{835519807167135744, 835519806932254720}

	service := article.ArticleService{}
	s, dbs := service.Find("6636feb4dec7395ca65806e87a5fe278", i, 1)
	idsByte, _ := json.Marshal(dbs)
	fmt.Println(s, string(idsByte))
}
