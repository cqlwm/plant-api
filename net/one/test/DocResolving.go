package test

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
)

func TestDoc() {
	getCityArr()
}

func getCityArr() {
	resp, err := http.Get(`http://127.0.0.1:8848/STX/the1.html`)
	if err != nil {
		fmt.Println("get request err", err)
		return
	}
	docByte := resp.Body
	doc, err := goquery.NewDocumentFromReader(docByte)
	if err != nil {
		log.Println(err)
		return
	}

	names := [100]string{}
	com := map[string]string{}
	doc.Find("#provinceSel option").Each(func(i int, selection *goquery.Selection) {
		names[i] = selection.Text()
		com[names[i]], _ = selection.Attr("value")
	})

	cps := make([]CityProvince, 0)
	for k, v := range com {
		list := getCityList(v)
		item := CityProvince{
			ProvinceName: k,
			Code:         v,
			City:         list,
		}
		cps = append(cps, item)
	}
	byss, _ := json.Marshal(cps)
	fmt.Println(string(byss))
}

func getCityList(code string) *[]CityEntry {
	resp, err := http.Get("http://www.nmc.cn/rest/province/" + code)
	if err != nil {
		panic("error get url")
	}
	res := []CityEntry{}
	bs, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(bs, &res)
	fmt.Println(res)
	return &res
}

type CityEntry struct {
	City string
	Code string
}

type CityProvince struct {
	ProvinceName string
	Code         string
	City         *[]CityEntry
}
