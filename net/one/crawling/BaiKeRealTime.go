package crawling

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
	"strings"
)

// 百度百科实时爬虫
func disease(name string) (*io.ReadCloser, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://baike.baidu.com/item/%s", name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求 ", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("响应错误 ", err)
		return nil, err
	}

	return &resp.Body, nil
}

// 寻找关键词
func Find(searchKey string) *entry.WeedResult {
	// 查找MongoDB
	weedResult := entry.WeedResult{}
	exist := config.CollectionGet(searchKey, &weedResult, config.WeedKey, config.WeedColl)
	if exist {
		log.Println("Find CollectionGet exist ", weedResult)
		return &weedResult
	}

	// 未找到，实时爬取
	body, err := disease(searchKey)
	if err != nil {
		log.Println("disease(searchKey) 未找到，实时爬取 ", err.Error())
		return &entry.WeedResult{}
	}

	doc, err := goquery.NewDocumentFromReader(*body)
	mainContent := doc.Find(".content-wrapper .content .main-content")
	// 名称
	title := mainContent.Find(".lemmaWgt-lemmaTitle-title h1").Text()

	// 目录
	catalog := make([]string, 0)
	mainContent.Find(".para-title h2").Each(func(i int, selection *goquery.Selection) {
		key := selection.Text()
		key = strings.TrimSpace(key)
		catalog = append(catalog, key[len(title):])
	})

	// 获取具体内容
	detailList := make([]string, 0)
	mainContent.Find("div").Each(func(i int, selection *goquery.Selection) {
		valClass, _ := selection.Attr("class")
		if valClass == "para-title level-2" {
			key := selection.Text()
			key = strings.TrimSpace(key)[len(title):]
			index1 := strings.Index(key, "\n") + 1
			key = strings.TrimSpace(key[:index1])
			detailList = append(detailList, key)
		}
		if valClass == "para" {
			key := selection.Text()
			key = strings.ReplaceAll(key, "\n", "")
			detailList = append(detailList, key)
		}
	})

	result := map[string]string{}
	result["searchKey"] = searchKey
	result["名称"] = title
	result["介绍"] = ""

	indexItem := 0
	for k, v := range detailList {
		if catalog[0] == v {
			indexItem = k
			break
		}
		result["介绍"] += v
	}

	for _, v := range catalog {
		result[v] = ""
	}

	itemKey := ""
	for loop2 := indexItem; loop2 < len(detailList); loop2++ {
		_, exist := result[detailList[loop2]]
		if exist {
			itemKey = detailList[loop2]
			continue
		}
		result[itemKey] += detailList[loop2]
	}

	weed := entry.Weeds{
		SearchKey: searchKey,
		Catalog:   catalog,
		Result:    result,
	}

	info, err := json.Marshal(weed)
	infoStr := string(info)
	log.Println("实时爬取到内容 1 ", infoStr)
	save := config.CollectionSave(searchKey, infoStr, config.WeedKey, config.WeedColl)

	if save {
		config.CollectionGet(searchKey, &weedResult, config.WeedKey, config.WeedColl)
	}

	log.Println("实时爬取到内容 2 ", weedResult)

	//fmt.Println("实时爬虫 ", weedResult)
	return &weedResult
}
