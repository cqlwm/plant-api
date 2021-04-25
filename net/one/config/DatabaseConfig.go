package config

import (
	_ "errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"plant-api/net/one/entry"
	"time"
)

var DataBase *gorm.DB

var autoMigrate = []interface{}{
	&entry.SysUser{}, &entry.Botany{}, &entry.PestTable{},
	&entry.IdentifyLog{}, &entry.ArticleDB{},
}

func Init() (bool, error) {
	dsn := "root:!BrWCrxPh0Vq@tcp(148.70.115.4:3306)/plant?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 数据库连接错误
	if err != nil {
		return false, NewError(err, &DATABASE_ERROR)
	}

	// 迁移 schema
	for _, v := range autoMigrate {
		db.AutoMigrate(v)
	}
	//db.AutoMigrate(&entry.SysUser{})
	//db.AutoMigrate(&entry.Botany{})
	//db.AutoMigrate(&entry.PestTable{})

	dbPool, err := db.DB()
	if err != nil {
		return false, NewError(err, &DATABASE_ERROR)
	}

	// 设置空闲连接池中连接的最大数量
	dbPool.SetMaxIdleConns(10)

	// 设置打开数据库连接的最大数量
	dbPool.SetMaxOpenConns(100)

	// 设置了连接可复用的最大时间
	dbPool.SetConnMaxLifetime(time.Hour)

	DataBase = db

	// 大部分 CRUD API 都是兼容的
	//db.AutoMigrate(&Product{})
	//db.Create(&user)
	//db.First(&user, 1)
	//db.Model(&user).Update("Age", 18)
	//db.Model(&user).Omit("Role").Updates(map[string]interface{}{"Name": "jinzhu", "Role": "admin"})
	//db.Delete(&user)
	return true, nil
}
