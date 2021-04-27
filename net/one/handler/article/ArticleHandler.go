package article

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"plant-api/net/one/config"
	"plant-api/net/one/service/impl/article"
	"plant-api/net/one/util"
	"strconv"
)

// handler路由配置
func ArticleHandler(e *gin.Engine) {
	e.POST(config.DistinguishHandler2search, search)
}

var service = article.ArticleService{}

// 搜索文章
func search(c *gin.Context) {
	// TODO 1.根据类型判断是图搜还是词搜
	// TODO 2.获取图片ids，相似度排序 前100张
	// TODO 3.生成分页ID组、SearchID、存入Redis
	// TODO 4.将第一个ID组作为结果集返回、并为前端生成分页信息

	// 默认图搜
	ts := false

	form, err := c.MultipartForm()
	if err != nil {
		// 表单提交异常
		return
	}

	var searchId = ""
	if len(form.Value["searchId"]) != 0 {
		searchId = form.Value["searchId"][0]
	}

	var page = 1
	if len(form.Value["page"]) != 0 {
		page, err = strconv.Atoi(form.Value["page"][0])
		if err != nil {
			config.Error2(c, "张林康页码用数字")
			return
		}
	}

	file := form.File["file"]
	if len(file) == 0 {
		ts = true
	}

	if ts {
		// 关键词搜索
		valueArr := form.Value["searchText"]
		if len(valueArr) == 0 {
			config.Error2(c, "参数有误")
			return
		}
		fmt.Println(valueArr[0])
		participle := util.SafeParticiple(valueArr[0])
		s, dbs := service.KeyFind(searchId, participle, page)

		resultOver := make(map[string]interface{})
		resultOver["searchId"] = s
		resultOver["article"] = dbs

		config.Ok2(c, resultOver)
		return
	}

	img := file[0]
	suffix, isImage := util.IsImage(img.Filename)
	if !isImage {
		config.Error2(c, config.PICTURE_FORMAT_ERROR.Message())
		return
	}

	// 保存到本地、调用识别
	completeNewName := util.Uuid(img.Filename) + suffix
	uploadSave := config.ImageSavePathConstant + completeNewName
	_ = c.SaveUploadedFile(img, uploadSave)

	ids, err := util.SearchByImage(uploadSave)

	// 接受文件或关键词
	s, dbs := service.Find(searchId, ids, page)

	resultOver := make(map[string]interface{})
	resultOver["searchId"] = s
	resultOver["article"] = dbs
	config.Ok2(c, resultOver)
}
