package logs

import (
	"log"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
)

type IdentifyLogServiceImpl struct{}

// 生成识别日志
func (ils *IdentifyLogServiceImpl) Log(logEntry *entry.IdentifyLog) {
	db := config.DataBase
	tx := db.Create(logEntry)
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Print("失败的日志", tx.Error.Error())
		return
	}
	if logEntry.Id != 0 {
		log.Println("成功的日志", *logEntry)
		return
	}
}

// 分页日志记录
func (ils *IdentifyLogServiceImpl) LogQuery(userId int, start int, size int) []entry.IdentifyLogSimple {
	db := config.DataBase
	sql := "select id, `option`, created, content  from identify_logs where user_id = ? limit ?, ?"
	//tx := db.Exec(sql, userId, start, size)
	rows, err := db.Debug().Raw(sql, userId, (start-1)*size, size).Rows()
	defer rows.Close()
	if err != nil {
		return nil
	}

	ilsArr := make([]entry.IdentifyLogSimple, 0)
	for rows.Next() {
		simple := entry.IdentifyLogSimple{}
		_ = db.ScanRows(rows, &simple)
		ilsArr = append(ilsArr, simple)
	}

	return ilsArr
}

// 总数
func (ils *IdentifyLogServiceImpl) LogCount(userId int) int {
	db := config.DataBase
	rows, _ := db.Raw("select count(user_id) as uNumber from identify_logs where user_id = ? ", userId).Rows()
	defer rows.Close()
	i := 0
	_ = rows.Next()
	rows.Scan(&i)
	return i
}
