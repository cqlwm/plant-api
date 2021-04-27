package entry

type IdentifyLog struct {
	Id       int `gorm:"primaryKey"`
	UserId   int
	Option   string
	OptionId int
	Content  string
	Created  int `gorm:"autoCreateTime"`
}

type IdentifyLogSimple struct {
	Id      int `gorm:"primaryKey"`
	Option  string
	Content string
	Created int `gorm:"autoCreateTime"`
}
