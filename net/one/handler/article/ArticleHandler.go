package article

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
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

type searchForm struct {
	// 1 图片；2 文字; 3 id
	pwd      int
	isErr    bool
	errInfo  config.CM
	searchId string
	page     int
	file     *multipart.FileHeader
	keyword  string
}

// 参数获取
func searchParam(c *gin.Context) *searchForm {
	resultForm := searchForm{}

	form, err := c.MultipartForm()
	if err != nil {
		resultForm.isErr = true
		resultForm.errInfo = &config.REQ_FORM_ERROR
		return &resultForm
	}

	resultForm.page = 1
	if len(form.Value["page"]) != 0 {
		page, err := strconv.Atoi(form.Value["page"][0])
		if err != nil {
			resultForm.isErr = true
			resultForm.errInfo = &config.PAGE_ERROR
			return &resultForm
		}
		if page <= 0 {
			page = 1
		}
		resultForm.page = page
	}

	if len(form.Value["searchId"]) != 0 {
		resultForm.searchId = form.Value["searchId"][0]
		if len(resultForm.searchId) > 0 {
			resultForm.pwd = 3
			return &resultForm
		}
	}

	file := form.File["file"]
	if len(file) != 0 {
		resultForm.pwd = 1
		resultForm.file = file[0]
		resultForm.pwd = 1
		return &resultForm
	}
	// 关键词搜索
	valueArr := form.Value["searchText"]
	if len(valueArr) == 0 {
		resultForm.isErr = true
		resultForm.errInfo = &config.Pricture_WORD_ERROR
		return &resultForm
	}
	resultForm.pwd = 2
	resultForm.keyword = valueArr[0]

	return &resultForm
}

// 词搜索
func searchWord(searchId string, searchText string, page int) map[string]interface{} {
	participle := util.SafeParticiple(searchText)
	s, dbs, forkPage := service.KeyFind(searchId, participle, page)

	resultOver := make(map[string]interface{})
	resultOver["searchId"] = s
	resultOver["article"] = dbs
	resultOver["page"] = forkPage

	return resultOver
}

// 图搜索
func searchPicture(c *gin.Context, searchId string, img *multipart.FileHeader, page int) {
	suffix, isImage := util.IsImage(img.Filename)
	if !isImage {
		config.Error(c, &config.PICTURE_FORMAT_ERROR, nil)
		return
	}

	// 保存到本地、调用识别
	completeNewName := util.Uuid(img.Filename) + suffix
	uploadSave := config.ImageSavePathConstant + completeNewName
	_ = c.SaveUploadedFile(img, uploadSave)

	ids, err := util.SearchByImage(uploadSave)
	if err != nil {
		config.Error(c, &config.SERVER_ERROR, nil)
		return
	}

	// 接受文件或关键词
	s, dbs, forkpage := service.FindPicture(ids, page)

	resultOver := make(map[string]interface{})
	resultOver["searchId"] = s
	resultOver["article"] = dbs
	resultOver["forkpage"] = forkpage
	config.Ok2(c, resultOver)
}

// 图搜索
func searchById(searchId string, page int) map[string]interface{} {
	dbs, forkPage := service.FindById(searchId, page)
	resultOver := make(map[string]interface{})
	resultOver["searchId"] = searchId
	resultOver["article"] = dbs
	resultOver["page"] = forkPage
	return resultOver
}

// 资讯搜索
func search(c *gin.Context) {
	param := searchParam(c)

	if param.isErr {
		config.Error(c, param.errInfo, nil)
		return
	}
	if param.pwd == 3 {
		seas := searchById(param.searchId, param.page)
		config.Ok2(c, seas)
		return
	}
	// 图
	if param.pwd == 1 {
		searchPicture(c, param.searchId, param.file, param.page)
		return
	}
	// 词
	if param.pwd == 2 {
		words := searchWord(param.searchId, param.keyword, param.page)
		config.Ok2(c, words)
		return
	}
	config.Error(c, &config.PARAM_ERROR, nil)
}

//
