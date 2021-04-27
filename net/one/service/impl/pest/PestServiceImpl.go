package pest

import (
	"fmt"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
)

// 根据害虫名称，获取害虫信息
func FindPestByName(name string) []entry.PestTable {
	pestArray := make([]entry.PestTable, 0)

	// 缓存
	key := fmt.Sprintf("%s%s", config.PestName, name)
	err := config.GetJson(key, &pestArray)
	if err == nil {
		fmt.Println("找到了...")
		return pestArray
	}

	sql := fmt.Sprintf(
		"select %s from pest_tables pt where name like ? or alias_name like ? or scientific_name like ? ",
		"id, name, alias_name, scientific_name, shape, habit, harm, parasitic, distribution, govern_method, pest_image")
	name = fmt.Sprintf("%%%s%%", name)

	db := config.DataBase
	rows, err := db.Debug().Raw(sql, name, name, name).Rows()
	defer rows.Close()

	if err != nil {
		return nil
	}

	for rows.Next() {
		table := entry.PestTable{}
		_ = db.ScanRows(rows, &table)
		pestArray = append(pestArray, table)
	}

	err = config.SaveKJson(key, pestArray, 0)
	if err != nil {
		fmt.Println(err)
	}

	return pestArray
}
