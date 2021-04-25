package entry

// 爬虫爬取文章的类
type Article struct {
	Id        int64
	Title     string
	Date      string
	Paragraph []string
	Original  string
}

// 对应数据库的类
type ArticleDB struct {
	Id        int64 `gorm:"primaryKey"`
	Title     string
	Date      string
	Paragraph string
	Original  string
	Created   int `gorm:"autoCreateTime"`
}
