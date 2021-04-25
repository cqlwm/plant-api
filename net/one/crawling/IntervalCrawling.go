package crawling

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
	"plant-api/net/one/util"
	"strings"
	"time"
)

var snowflakeIdWorker = util.SnowflakeIdWorker{}

// 定时爬虫
// 定时爬取信息到数据库
// https://www.my478.com/html/list/ask/
// http://www.86362.com/
// TODO 待添加定时器
func IntervalAskUrl() {
	url := `https://www.my478.com/html/list/ask/`
	resp, err := http.Get(url)
	if err != nil {
		log.Println("IntervalAskUrl 异常请求...")
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("IntervalAskUrl 响应体解析异常")
		return
	}

	now := date()
	doc.Find("body > div.web > div.list-r.list-z > ul").
		Each(func(i int, selection *goquery.Selection) {
			hrefA, _ := selection.Find("a").Attr("href")
			span := selection.Find("span").Text()
			if strings.Index(span, now) != -1 {
				go detailAskUrl("https://www.my478.com/" + hrefA)
			}
		})
}

// 爬取详细页
func detailAskUrl(detail string) {
	resp, err := http.Get(detail)
	if err != nil {
		log.Println("Detail 异常请求...")
		return
	}
	defer resp.Body.Close()

	// doc
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("DetailAskUrl 文档解析异常")
		return
	}

	content := doc.Find("body > div.web > div.moban")
	title := content.Find("h1").Text()
	date := content.Find(".date").Text()

	id := snowflakeIdWorker.NextId()
	article := entry.ArticleDB{
		Id:    id,
		Title: title,
		Date:  date,
	}

	paragraph := make([]string, 0)
	intro := content.Find(".intro").Text()
	paragraph = append(paragraph, intro)
	// 图片
	originSrcImage, _ := content.Find(".center > img").Attr("src")
	imgSuffix, ok := util.IsImage(originSrcImage)
	newFileName := ""
	if ok {
		down := `https://www.my478.com/` + originSrcImage
		newFileName = fmt.Sprintf("%d%s", article.Id, imgSuffix)
		save := config.SearchByImage + newFileName
		util.DownImage2(save, down)
		paragraph = append(paragraph, config.ImageTag+config.SearchByImageHttpUrl+newFileName)
	}
	// 正文
	content.Find(".content > p").Each(func(i int, selection *goquery.Selection) {
		paragraph = append(paragraph, selection.Text())
	})
	pJson, _ := json.Marshal(paragraph)
	article.Paragraph = string(pJson)
	article.Original = detail

	b, _ := json.Marshal(article)
	//fmt.Println(string(b))

	db := config.DataBase
	tx := db.Save(&article)
	i := tx.RowsAffected
	fmt.Println(i, string(b))
}

// date
func date() string {
	year, month, day := time.Now().Date()
	nowString := fmt.Sprintf("%d-%02d-%d", year, month, day)
	return nowString
}
