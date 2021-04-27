package entry

type DistinguishQuery struct {
	Choice        string
	TargetPicture string
}

/*
	Current  请求页
	PageSize 页大小
*/
type HistoryForm struct {
	Current  int `json:"Current"  binding:"required"`
	PageSize int `json:"PageSize" binding:"required"`
}
