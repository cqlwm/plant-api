package config

import "runtime"

var finalBase1 = `/root/plant/`
var finalBase2 = `/usr/local/nginx/html/plant/`
var urlBase = `https://graceful.top/`

// 城市JSON文件
var CityJsonPath = `/root/plant/staticfile/city.json`

// 分词词库路径
var SegmenterLoadDictionary = finalBase1 + "staticfile/dictionary.txt"

// 临时图片文件：本地、网络地址
var ImageSavePathConstant = finalBase2 + `image/item/`
var ReUrlImage = urlBase + `plant/image/item/`

// 爬取每日最新资讯，图片保存路径；对应以图搜图的功能；
var SearchByImage = finalBase2 + `image/searchby/`
var SearchByImageHttpUrl = urlBase + `plant/image/searchby/`

// 端口号
var Port = ":9091"

// ######################### MongoDB Key #########################
// 杂草识别模块
var WeedKey = "WeedKey"
var WeedColl = "weeds"

// 定时文章爬虫，存入数据库时的Tag
var ImageTag = "ThisImageTag:"

// ######################### Redis Key #########################
// 文章搜索IDS
var ArticleIds = "ArticleSearchId:"

// 天气Key
var WeatherRedisKey = "WeatherRedisKey-code:"
var SimpleWeatherRedisKey = "SimpleWeatherRedisKey-code:"

// Windows环境加载
func OsLoad() {
	osName := runtime.GOOS
	winPath := `F:\deskFile\02智慧农业\plant`
	if osName == "windows" {
		// 临时图片放置
		ImageSavePathConstant = winPath + `\test\`
		// 搜索图片
		SearchByImage = winPath + `\searchby\`
		// 分词词库路径
		SegmenterLoadDictionary = winPath + `\dictionary.txt`
		// 城市JSON文件
		CityJsonPath = `net/one/config/city.json`
	}
}
