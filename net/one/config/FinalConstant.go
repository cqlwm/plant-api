package config

import "runtime"

// 分词词库路径
var SegmenterLoadDictionary = "/root/plant/staticfile/dictionary.txt"

// 临时图片文件：本地、网络地址
var ImageSavePathConstant = `/usr/local/nginx/html/plant/image/item/`
var ReUrlImage = `https://graceful.top/plant/image/item/`

// 爬取每日最新资讯，图片保存路径；对应以图搜图的功能；
var SearchByImage = `/usr/local/nginx/html/plant/image/searchby/`
var SearchByImageHttpUrl = `https://graceful.top/plant/image/searchby/`

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
	}
}
