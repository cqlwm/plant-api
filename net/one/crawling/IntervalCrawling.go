package crawling

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"plant-api/net/one/config"
	"plant-api/net/one/util"
)

// 定时爬虫
// 定时爬取信息到数据库
// https://www.my478.com/html/list/ask/
// http://www.86362.com/
func IntervalAskUrl() {
	url := `https://www.my478.com/html/list/ask/`
	resp, err := http.Get(url)
	if err != nil {
		log.Println("IntervalAskUrl 异常请求...")
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("IntervalAskUrl 响应体解析异常")
		return
	}

	doc.Find("body > div.web > div.list-r.list-z > ul").
		Each(func(i int, selection *goquery.Selection) {
			nodeA := selection.Find("a")
			href, _ := nodeA.Attr("href")
			text := nodeA.Text()

			span := selection.Find("span").Text()
			fmt.Println(href, text, span)
		})

	//详细页
	// https://www.my478.com/html/20210423/373659.html
	// /html/20210423/373659.html
}

// 爬取详细页
func DetailAskUrl(detail string) {
	resp, err := http.Get(detail)
	if err != nil {
		log.Println("Detail 异常请求...")
		return
	}

	// doc
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("DetailAskUrl 文档解析异常")
		return
	}

	content := doc.Find("body > div.web > div.moban")
	title := content.Find("h1").Text()
	date := content.Find(".date").Text()
	intro := content.Find(".intro").Text()
	originSrcImage, _ := content.Find(".center > img").Attr("src")
	contentArr := make([]string, 0)
	content.Find(".content > p").Each(func(i int, selection *goquery.Selection) {
		contentArr = append(contentArr, selection.Text())
	})

	fmt.Println(title)
	fmt.Println(date)
	fmt.Println(intro)

	// 保存在本地
	imgSuffix, ok := util.IsImage(originSrcImage)
	// TODO 不适用UUID，使用雪花ID代替图片名称
	newFileName := ""
	if ok {
		down := `https://www.my478.com/` + originSrcImage
		newFileName = util.Uuid(originSrcImage) + imgSuffix
		save := config.SearchByImage + newFileName
		fmt.Printf("down: %s, save: %s\n", down, save)
		util.DownImage(save, down)
	}

	for k, v := range contentArr {
		fmt.Println(k, v)
	}
}
