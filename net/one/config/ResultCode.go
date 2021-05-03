package config

import (
	"errors"
	"fmt"
)

// ##########200_000 成功##########
var SUCCESS = codeModel{200_000, "成功"}

// ##########300_000 参数与权限##########
var PARAM_ERROR = codeModel{300_000, "参数错误"}
var FORM_ERROR = codeModel{300_001, "请以表单上传，或表单参数错误"}
var PICTURE_PARAM_ERROR = codeModel{300_002, "图片上传失败"}
var PICTURE_FORMAT_ERROR = codeModel{300_003, "不正确的图片格式"}
var TOKEN_INVALID_ERROR = codeModel{300_004, "Token Invalid or expired"}
var WX_CODE_ERROR = codeModel{300_005, "Token Invalid or expired"}
var NIL_TOKEN_ERROR = codeModel{300_006, "空的Token参数"}
var NOMORE_DATA_ERROR = codeModel{300_007, "已经没有更多的数据了"}
var REQ_FORM_ERROR = codeModel{300_008, "请使用表单提交"}
var PAGE_ERROR = codeModel{300_009, "页码请使用整数"}
var Pricture_WORD_ERROR = codeModel{300_010, "图搜词搜必选其一"}

// ##########400_000 前端异常##########

// ##########500_000 服务器异常##########
var SERVER_ERROR = codeModel{500_000, "服务器繁忙"}
var DATABASE_ERROR = codeModel{500_001, "数据库连接异常"}
var JSON_SERIALIZE_ERROR = codeModel{500_002, "JSON序列化异常"}

// ##########600_000 逻辑异常##########
var PICTURE_ANALYSIS_ERROR = codeModel{600_001, "图片解析异常"}

// 响应统一接口
type CM interface {
	State() int
	Message() string
}

// 响应私有实现模型
type codeModel struct {
	state   int
	message string
}

func (cm *codeModel) State() int {
	return cm.state
}

func (cm *codeModel) Message() string {
	return cm.message
}

func NewError(e error, cm *codeModel) error {
	_ = fmt.Errorf(e.Error())
	err := fmt.Sprintf("%d:%s", cm.state, cm.message)
	return errors.New(err)
}
