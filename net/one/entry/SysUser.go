package entry

// 数据库实体对象
type SysUser struct {
	Id int `gorm:"primaryKey"`
	// 随机id
	UUID string
	// 用户名
	UserName string
	// WxID
	OpenId string
	// 昵称
	NickName string
	// 头像
	AvatarUrl string
	// 性别
	Gender int
	// 国家
	Country string
	// 省
	Province string
	// 城市
	City string
}

// 登录表单
type LoginForm struct {
	Code      string `form:"code" json:"code" binding:"required"`
	AvatarUrl string `form:"avatarUrl" json:"avatarUrl" binding:"required"`
	City      string
	Country   string
	Gender    int
	Language  string
	NickName  string `form:"nickName" json:"nickName" binding:"required"`
	Province  string
}
