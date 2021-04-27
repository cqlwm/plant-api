package tool

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
	"plant-api/net/one/service/impl/logs"
	"plant-api/net/one/service/impl/pest"
	"plant-api/net/one/util"
	"strings"
)

var savePath = config.ImageSavePathConstant

// handler路由配置
func DistinguishHandler(e *gin.Engine) {
	e.POST(config.DistinguishHandlerQuery, query)
	e.POST(config.DistinguishHandler2HistoryRecording, HistoryRecording)
}

// 参数
func getQueryParam(c *gin.Context) (*multipart.FileHeader, string, config.CM) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, "", &config.FORM_ERROR
	}
	// 获取图片
	pictures := form.File["picture"]
	if pictures == nil {
		return nil, "", &config.PICTURE_PARAM_ERROR
	}
	_, isImage := util.IsImage(pictures[0].Filename)
	if !isImage {
		return nil, "", &config.PICTURE_FORMAT_ERROR
	}
	// 获取选择参数
	choice := form.Value["Choice"]
	if choice == nil || len(choice) == 0 || !isFourOne(choice[0]) {
		return nil, "", &config.PARAM_ERROR
	}
	return pictures[0], choice[0], nil
}

// 判断模式中的四个之一
func isFourOne(choice string) bool {
	if len(choice) == 0 {
		return false
	}
	choice = strings.ToLower(choice)
	switch choice {
	case "pest", "weeds", "pesticide", "estimate", "disease":
		return true
	}
	return false
}

// 识别接口
func query(c *gin.Context) {
	claims, errParm := util.UserInfo(c)
	if errParm != nil {
		config.Error2(c, errParm.Error())
		return
	}

	file, choice, ret := getQueryParam(c)

	if ret != nil {
		config.Error(c, ret, nil)
		return
	}

	// 保存图片 TODO 图片识别之后生成日志
	completeNewName := util.Uuid(file.Filename) + util.Extension(file.Filename)
	uploadSave := savePath + completeNewName
	_ = c.SaveUploadedFile(file, uploadSave)

	// 害虫识别
	ds := pest.DistinguishServiceImpl{}
	//var resArr []entry.PestSimilarity

	var err error
	var obj interface{}
	// "pest", "weeds", "pesticide", "estimate":  根据不同模式调用不同识别模块
	switch choice {
	case "pest":
		obj, err = ds.Pest(uploadSave)
		//obj = resArr
	case "weeds":
		obj, err = nil, nil
	case "disease":
		disease := ds.Disease(uploadSave, completeNewName)
		//ret = disease
		if disease == nil {
			config.Error2(c, nil)
			return
		}
		obj = disease
		return
	case "pesticide":
		obj, err = nil, nil
	case "estimate":
		obj, err = nil, nil
	}

	if err != nil {
		config.Error2(c, nil)
		return
	}

	objByte, err := json.Marshal(obj)

	// 生成操作日志
	ilog := entry.IdentifyLog{
		UserId:  claims.UserId,
		Option:  choice,
		Content: string(objByte),
	}
	logService := logs.IdentifyLogServiceImpl{}
	logService.Log(&ilog)

	// 结果集
	config.Ok2(c, obj)
}

// 识别记录
func HistoryRecording(c *gin.Context) {
	//claims, err := util.UserInfo(c)
	//if err != nil {
	//	config.Ok2(c, err.Error())
	//	return
	//}

	claims := util.CustomClaims{
		UserId: 100,
	}

	form := entry.HistoryForm{}
	_ = c.ShouldBindJSON(&form)

	ils := logs.IdentifyLogServiceImpl{}
	// 验证页码是否溢出
	count := ils.LogCount(claims.UserId)
	page := entry.ForkPage{}
	page.Set(form.Current, form.PageSize, count)
	if page.Overflow() {
		page.Content = []entry.IdentifyLogSimple{}
		config.Ok(c, &config.NOMORE_DATA_ERROR, page)
		return
	}
	// 查找页码
	logQuery := ils.LogQuery(claims.UserId, form.Current, form.PageSize)
	page.Content = logQuery
	config.Ok2(c, page)

}
