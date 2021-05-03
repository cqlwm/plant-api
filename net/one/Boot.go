package main

import (
	"fmt"
	"plant-api/net/one/config"
	"plant-api/net/one/handler"
	"plant-api/net/one/handler/article"
	"plant-api/net/one/handler/tool"
	"plant-api/net/one/handler/user"
	"plant-api/net/one/handler/weather"
)

func main() {
	// 参数加载
	config.OsLoad()

	// GRom
	ok, err := config.Init()
	if err != nil {
		fmt.Println("GRom 数据源 : ", err)
		return
	}
	if ok {
		fmt.Println("数据库初始化完成")
	}

	// Redis
	err = config.InitRedisClient()
	if err != nil {
		fmt.Println("Redis 初始化 : ", err)
		return
	}

	// MongoDB
	config.MongoDBConn()

	// Gin
	r := config.R()
	handler.TestHandlerRoute(r) // 测试
	tool.DistinguishHandler(r)  // 虫、草、药识别
	user.UserHandler(r)         // 用户
	article.ArticleHandler(r)   // 文章
	weather.WeatherHandler(r)   // 天气

	r.Run(config.Port)
}

func main2() {
	//crawling.PestWebRequest()
	//srcByte, err := ioutil.ReadFile(`C:\Users\lwm\Desktop\1\1.jpg`)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//res := base64.StdEncoding.EncodeToString(srcByte)
	//
	//fmt.Println(res)
}
