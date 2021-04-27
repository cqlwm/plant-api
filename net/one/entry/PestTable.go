package entry

// 害虫表
type PestTable struct {
	// 主键ID
	Id int `gorm:"primaryKey"`
	// 名称
	Name string
	// 别名
	AliasName string
	// 学名
	ScientificName string
	// 形态
	Shape string
	// 习性
	Habit string
	// 特点
	Harm string
	// 寄主
	Parasitic string
	// 分布
	Distribution string
	// 防治方法
	GovernMethod string
	// 图像
	PestImage string
}
