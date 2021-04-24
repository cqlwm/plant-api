package config

const (
	//ReUrlImage = `http://127.0.0.1:8848/browser/image/`
	//ImageSavePathConstant = `C:\Users\lwm\Desktop\demo-for-aliyun-plant-api-master\javascript\browser\image\`

	SearchByImage = `F:\test\images\plant\searchby\`

	// 生产环境
	ImageSavePathConstant = `/usr/local/nginx/html/plant/image/item/`
	ReUrlImage            = `http://graceful.top/plant/image/item/`
	// 爬取每日最新资讯，图片保存路径；对应以图搜图的功能；
	// SearchByImage = `/usr/local/nginx/html/plant/image/searchby/`

	Port = ":9091"

	// MongoDB Key
	WeedKey  = "WeedKey"
	WeedColl = "weeds"
)
