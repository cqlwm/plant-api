package crawling

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"plant-api/net/one/config"
	"plant-api/net/one/entry"
	"plant-api/net/one/util"
	"strings"
	"time"
)

// http://www.51agritech.com/pestweb/sc/SCHC.htm
func PestWebRequest() {
	baseUrl := "http://www.51agritech.com/pestweb/sc/SCHC.htm"
	result := Get(baseUrl)

	fmt.Println(result)
}

func Base() {
	baseUrl := "http://www.51agritech.com/pestweb/sc/SCHC.htm"
	res, err := http.Get(baseUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()
	fmt.Printf(newStr)

	//body, err := iconv.ConvertString(newStr, "GBK", "utf-8")
	//fmt.Println(body)
}

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(url string) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	d, err := util.GbkToUtf8(result.Bytes())
	return string(d)
}

// 爬取列表
func ExampleScrape() {
	// Request the HTML page.http://www.51agritech.com/pestweb/sc/SCHC.htm
	res, err := http.Get("http://www.51agritech.com/pestweb/sc/SCHC.htm")
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	val, exists := doc.Find("form > input").Attr("value")
	if !exists {
		fmt.Println("没有找到这个属性")
	}
	r, _ := util.GbkToUtf8([]byte(val))
	resultOrg := string(r)

	s1 := strings.ReplaceAll(resultOrg, "|", "")
	s2 := strings.Split(s1, "*.")

	base := "http://www.51agritech.com/pestweb/sc/"
	for _, v := range s2 {
		v = strings.ReplaceAll(v, "./V  HC", "")
		v = strings.ReplaceAll(v, "~", "")
		indexLast := strings.LastIndex(v, ".htm")
		v := v[:indexLast+4]

		DetailedImage(base + v)

		detailed, err := Detailed(base + v)
		if err == nil {
			//detailed.Id = k
			//fmt.Println(k)
			config.DataBase.Save(detailed)
		}
	}
}

func DetailedImage(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("status code error:", res.StatusCode, res.Status)
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	val, _ := doc.Find("img").Attr("src")
	index := strings.LastIndex(url, "/")
	url = url[:index+1] + val

	//	// 保存图片
	save := `C:\Users\lwm\Desktop\1\The\`
	util.DownImage(save, url)
	return val
}

// 爬取详细内容
func Detailed(url string) (*entry.PestTable, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("status code error:", res.StatusCode, res.Status)
		return nil, errors.New("status code error")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	pest := entry.PestTable{}

	claes := []string{".MsoNormal", ".MsoPlainText"}
	for ci := 0; ci < len(claes); ci++ {
		doc.Find(claes[ci]).Each(func(i int, s *goquery.Selection) {
			item := strings.ReplaceAll(strings.TrimSpace(s.Text()), "\n", "")

			alias := "别名"
			if index := strings.Index(item, alias); index == 0 {
				item = strings.Replace(item, alias, "", len(alias))
				item = strings.TrimSpace(item)
				pest.AliasName = item
			}

			scientificName := "学名"
			if index := strings.Index(item, scientificName); index == 0 {
				item = strings.Replace(item, scientificName, "", len(scientificName))
				item = strings.TrimSpace(item)
				pest.ScientificName = item
			}

			shape := "形态特征"
			if index := strings.Index(item, shape); index == 0 {
				item = strings.Replace(item, shape, "", len(shape))
				item = strings.TrimSpace(item)
				pest.Shape = item
			}

			habit := "生活习性" // 发生规律
			hasHabit := strings.Index(item, habit)
			if hasHabit == -1 {
				habit = "发生规律"
				hasHabit = strings.Index(item, habit)
			}
			if hasHabit == 0 {
				item = strings.Replace(item, habit, "", len(habit))
				item = strings.TrimSpace(item)
				pest.Habit = item
			}

			harm := "为害特征"
			hasHarm := strings.Index(item, harm)
			if hasHarm == -1 {
				harm = "为害特点"
				hasHarm = strings.Index(item, harm)
			}
			if hasHarm == 0 {
				item = strings.Replace(item, harm, "", len(harm))
				item = strings.TrimSpace(item)
				pest.Harm = item
			}

			parasitic := "寄主"
			if index := strings.Index(item, parasitic); index == 0 {
				item = strings.Replace(item, parasitic, "", len(parasitic))
				item = strings.TrimSpace(item)
				pest.Parasitic = item
			}

			distribution := "分布"
			if index := strings.Index(item, distribution); index == 0 {
				item = strings.Replace(item, distribution, "", len(distribution))
				item = strings.TrimSpace(item)
				pest.Distribution = item
			}

			// 防治方法
			govern := "防治方法"
			if index := strings.Index(item, govern); index == 0 {
				item = strings.Replace(item, govern, "", len(govern))
				item = strings.TrimSpace(item)
				pest.GovernMethod = item
			}

			if i == 0 && len(item) <= 45 {
				pest.Name = item
			}

		})
		if len(pest.Name) != 0 {
			break
		}
	}

	// 防止治理方法获取不到
	if len(pest.GovernMethod) == 0 {
		governText := doc.Find("body > span").Text()
		governText = strings.ReplaceAll(strings.TrimSpace(governText), "\n", "")
		governText = strings.Replace(governText, "防治方法", "", len("防治方法"))
		governText = strings.TrimSpace(governText)
		pest.GovernMethod = governText
	}

	// 保存图片
	pest.PestImage = DetailedImage(url)

	return &pest, nil
}
