package article

import (
	"fmt"
	"log"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
	"plant-api/net/one/util"
	"sort"
	"strings"
	"time"
)

type ArticleService struct{}

var pageSize = 15

// 图搜
func (as *ArticleService) Find(searchID string, ids []int, page int) (string, []entry.ArticleDB) {
	if len(searchID) != 0 {
		findIds := make([]entry.ArticleDB, 0)
		err := config.GetJson(config.ArticleIds+searchID, &findIds)
		if err == nil {
			return searchID, findIds
		}
		fmt.Println(err)
	}

	byIds := selectArticleByIds(ids, page)
	searchID = buildSearchID(byIds)

	return searchID, byIds
}

func (as *ArticleService) KeyFind(searchID string, keywords []string, page int) (string, []entry.ArticleDB) {

	word := make([]entry.ArticleDB, 0)
	if len(searchID) != 0 {
		err := config.GetJson(config.ArticleIds+searchID, &word)
		fmt.Println(err)
	}

	if len(word) == 0 {
		word = selectByKeyWord(keywords)
	}

	forkPage := entry.ForkPage{}
	forkPage.Set(page, pageSize, len(word))

	uuid := buildSearchID(word)

	var start = (page - 1) * pageSize
	var end = start + pageSize

	if start >= forkPage.Total {
		start = forkPage.Total - 1
	}

	if end > forkPage.Total {
		end = forkPage.Total
	}

	return uuid, word[start:end]

}

type artList []articleSuitability

func (I artList) Len() int {
	return len(I)
}
func (I artList) Less(i, j int) bool {
	return I[i].Sui > I[j].Sui
}
func (I artList) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

type articleSuitability struct {
	Sui     int
	Article entry.ArticleDB
}

// 数据库查找
func selectArticleByIds(ids []int, page int) []entry.ArticleDB {
	forkPage := entry.ForkPage{}
	forkPage.Set(page, pageSize, len(ids))

	startIndex := (page - 1) * pageSize

	id2match := map[int]int{}
	arrArticleDB := make([]entry.ArticleDB, 0)

	sql1 := `select id, title, date, paragraph, original, created from article_dbs `
	sql2 := `where id in `
	sql3 := "("

	end := startIndex + pageSize
	if end > len(ids) {
		end = len(ids)
	}
	for k, v := range ids[startIndex:end] {
		sql3 = fmt.Sprintf("%s%v,", sql3, v)
		id2match[v] = k
		arrArticleDB = append(arrArticleDB, entry.ArticleDB{})
	}
	sql3 = fmt.Sprintf("%s%s", sql3[:len(sql3)-1], ")")
	sql := fmt.Sprintf("%s%s%s", sql1, sql2, sql3)

	db := config.DataBase
	rows, err := db.Debug().Raw(sql).Rows()
	if err != nil {
		log.Println("SQL执行异常...")
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		article := entry.ArticleDB{}
		err = db.ScanRows(rows, &article)
		if err == nil {
			indexKey := id2match[int(article.Id)]
			arrArticleDB[indexKey] = article
		}
	}
	return arrArticleDB
}

// 关键词分页查找
func selectByKeyWord(keywords []string) []entry.ArticleDB {
	sql1 := "SELECT id, title, date, paragraph, original, created from article_dbs WHERE "
	sql2 := ""
	for k, v := range keywords {
		sor := "or"
		if k == len(keywords)-1 {
			sor = ""
		}
		sql2 += fmt.Sprintf(" paragraph like '%%%s%%' %s ", v, sor)
	}
	sql := fmt.Sprintf("%s%s limit 0, 75", sql1, sql2)

	db := config.DataBase
	rows, err := db.Debug().Raw(sql).Rows()
	if err != nil {
		// SQL 错误
		return nil
	}
	defer rows.Close()

	arDbRest := make([]articleSuitability, 0)
	for rows.Next() {
		arDb := entry.ArticleDB{}
		_ = db.ScanRows(rows, &arDb)
		count := 0
		for _, v := range keywords {
			if strings.Index(arDb.Paragraph, v) != -1 {
				count++
			}
		}
		arDbRest = append(arDbRest, articleSuitability{
			Sui:     count,
			Article: arDb,
		})
	}
	sort.Sort(artList(arDbRest))

	dbs := make([]entry.ArticleDB, 0)

	for _, v := range arDbRest {
		dbs = append(dbs, v.Article)
	}

	return dbs
}

// 生成searchID、排序、Redis 返回第一个
func buildSearchID(adb interface{}) string {
	uuid := util.Uuid(string(time.Now().Nanosecond()))
	err := config.SaveKJson(config.ArticleIds+uuid, adb, 5*60)
	if err != nil {
		log.Println(err)
	}
	return uuid
}
