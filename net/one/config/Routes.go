package config

const (
	base = "/plant-api"

	// ######################
	TestHandlerPing = base + "/ping"

	// ######################
	UserHandlerBase   = base + "/user/user"
	UserHandler2Login = UserHandlerBase + "/login"

	// ######################
	DistinguishHandlerBase              = base + "/tool/distinguish"
	DistinguishHandlerQuery             = DistinguishHandlerBase + "/query"
	DistinguishHandler2HistoryRecording = DistinguishHandlerBase + "/HistoryRecording"
)
