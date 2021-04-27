package entry

type WeedFlower struct {
	// 姓名
	Name string
	// 相似度
	Score string
	// 属名
	Genus string
	// 科名
	Family string
	// 图片
	ImageUrl string
	// Code码
	InfoCode string
}
type WeedFlowerResult struct {
	Status  int
	Message string
	Result  []WeedFlower
}

type Weeds struct {
	// 名称
	SearchKey string
	// 目录
	Catalog []string
	// 内容
	Result map[string]string
	//// 简介
	//Introduction string
	//// 发病原因
	//Cause string
	//// 发病症状
	//OnsetSymptoms string
	//// 发病条件
	//OnsetCondition string
	//// 防止方法
	//PreventionMethod string
}

type WeedResult struct {
	WeedKey string
	Info    string
}
