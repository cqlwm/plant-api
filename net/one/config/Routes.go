package config

const (
	base = "/plant-api"

	// ###################### 测试
	TestHandlerPing = base + "/ping"

	// ###################### 用户
	UserHandlerBase   = base + "/user/user"
	UserHandler2Login = UserHandlerBase + "/login"

	// ###################### 识别
	DistinguishHandlerBase              = base + "/tool/distinguish"
	DistinguishHandlerQuery             = DistinguishHandlerBase + "/query"
	DistinguishHandler2HistoryRecording = DistinguishHandlerBase + "/HistoryRecording"

	// ###################### 文章搜搜
	ArticleHandlerBase        = base + "/article"
	DistinguishHandler2search = ArticleHandlerBase + "/search"
)
