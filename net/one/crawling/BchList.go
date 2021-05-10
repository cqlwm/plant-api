package crawling

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

// 病虫害分页列表
var bchListUrl = `http://zhibao.yuanlin.com/bchList.aspx?page=`

// 病虫害详细信息
var bchDetailUrl = `http://zhibao.yuanlin.com/bchDetail.aspx?ID=`

func BchInfo(id string) {
	resp, err := getResp(id)
	if err != nil {
		return // Todo
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	bch := BchEntry{}

	// Image1
	bch.Name = doc.Find("#lb_mingc").Text()
	bch.PestImage = doc.Find("#Image1").AttrOr("src", "")
	bch.ScientificName = doc.Find("#lb_xueming").Text()
	bch.Generic = doc.Find("#lb_leishu").Text()

	pics := make([]string, 0)
	doc.Find("body .newslist .box_content2 .bch_text.tcenter img").
		Each(func(i int, selection *goquery.Selection) {
			src := selection.AttrOr("src", "")
			pics = append(pics, src)
		})
	bch.RelatedPictures = pics

	text := doc.Find(".newslist .box_content2 .bch_cont .bch_c2 .bch_text").Text()
	bch.DistributionAndHazards = trims(text)

	//log.Println(trims(text))
}

func BchList(page string) {

}

func getResp(s string) (*http.Response, error) {
	partialUrl := fmt.Sprintf("%s%s", bchDetailUrl, s)
	resp, err := http.Get(partialUrl)
	return resp, err
}

func trims(s string) string {
	space := strings.TrimSpace(s)
	return space
}

type BchEntry struct {
	Id int `gorm:"primaryKey"`
	// 名称
	Name string
	// 学名
	ScientificName string
	// 类属
	Generic string
	// 图像
	PestImage string
	// 相关图片
	RelatedPictures []string
	// 分布与危害
	DistributionAndHazards string
	// 形态特征
	Shape string
	// 发生规律-习性
	Habit string
	// 防止方法
	GovernMethod string
	// 防止药械
	GovernMedicine string
}
