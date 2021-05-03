package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/types"
	"gopkg.in/fatih/set.v0"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"plant-api/net/one/config"
	"plant-api/net/one/crawling"
	"plant-api/net/one/entry"
	"plant-api/net/one/util"
	"reflect"
	"strings"
)

func main00() {
	//participle := util.SafeParticiple("板栗果仁腐烂是什么病害，怎么防治？")
	//for k, v := range participle {
	//	fmt.Println(k, v, len(v))
	//}
	//week, _ := crawling.FutureWeek(57512)
	//b, _ := json.Marshal(week)
	//fmt.Println(string(b))

}

/*
什么
原因
*/

type TestSet struct {
	Id   int
	Name string
}

func setTest() {
	g := set.New(set.ThreadSafe)
	g.Add(TestSet{
		Id:   1,
		Name: "tom",
	})

	a := set.New(set.NonThreadSafe)
	a.Add(TestSet{
		Id:   1,
		Name: "tom",
	})
	a.Add(TestSet{
		Id:   2,
		Name: "hello",
	})
	a.Add(TestSet{
		Id:   3,
		Name: "world",
	})

	b := set.New(set.ThreadSafe)
	b.Add(TestSet{
		Id:   2,
		Name: "hello",
	})
	b.Add(TestSet{
		Id:   3,
		Name: "world",
	})
	b.Add(TestSet{
		Id:   4,
		Name: "jack",
	})

	//并集
	unionSet := set.Union(a, b, g)
	fmt.Printf("并集:%v\n", unionSet)

	//交集
	intersectionSet := set.Intersection(a, b, g)
	fmt.Printf("交集:%v\n", intersectionSet)

	//差集
	diffS1S2 := set.Difference(a, b)
	fmt.Printf("差集(属a不属b):%v\n", diffS1S2)

	diffS2S1 := set.Difference(b, a)
	fmt.Printf("差集(属b不属a):%v\n", diffS2S1)
}

func arrTest(i []int) {
	i[0] = 100
}

type A struct {
	Name string
	Age  int
}
type B struct {
	Name string
}

func jsonTest() {
	a := A{
		Name: "12345",
		Age:  10,
	}
	b := B{}
	to := util.BeanTo(a, &b)
	fmt.Println(to, b)
}

func wxTest() {
	id, key, ok := util.WeiXinAuth("033rYb0002DuyL1Xry100INIdV1rYb0a")
	fmt.Println(id, key, ok)
}

func nongye() {
	/*
		识别类型 1:柑橘 2:杂草 123:苹果 127:葡萄 131:小麦 132:水稻
		142:草地贪夜蛾 143:虫体
	*/
	res, err := util.PostFlower(`C:\Users\lwm\Desktop\1\GJ01.jpg`)
	//
	//// 辨认识别物
	key := ""
	if err == nil {
		// 识别
		typeMap := map[string]string{
			"柑橘": "1",
			"柑":  "1",
			"橘":  "1",
			"橙":  "1",
			"苹果": "123",
			"葡萄": "127",
			"小麦": "131",
			"麦":  "131",
			"水稻": "132",
			"稻":  "132",
		}
		for _, v2 := range res.Result {
			for k, v1 := range typeMap {
				index := strings.Index(v2.Name, k)
				if index != -1 {
					key = v1
					break
				}
			}
			if key != "" {
				break
			}
		}
		fmt.Println(key)
	}
	if key == "" {
		key = "2"
	}
	fmt.Println(key)

	//fmt.Println(typeMap)
	dire, err := util.Disease(key, `http://graceful.top/plant/image/other/GJ01.jpg`)
	if err != nil {
		fmt.Println(err)
	}
	diseaseName := dire.Content.Result

	// 百度百科爬虫
	config.MongoDBConn()
	config.DeleteAll(config.WeedColl)
	find := crawling.Find(diseaseName)
	fmt.Println(*find)

}

// MongoDB存取测试
func mongodbTest() {
	config.MongoDBConn()
	//crawling.Find("立枯病")
	//catalog := []string{"hello"}
	//result := map[string]string{
	//	"hello": "world",
	//}
	//e := entry.Weeds{
	//	SearchKey: "黑痘病Test",
	//	Catalog:   catalog,
	//	Result:    result,
	//}
	//bytes, _ := json.Marshal(e)
	//
	//save := config.CollectionSave("Test001", string(bytes), "weeds")
	//fmt.Println(save)

	// ObjectID("607bf260a7bfe034b93bb7a6")
	data := entry.WeedResult{}
	config.CollectionGet("Test001", &data, config.WeedKey, config.WeedColl)
	fmt.Println(data)
}

func urlCode() {
	escape := url.QueryEscape("//Hello/")
	fmt.Println(escape)
}

func sqlSubRep() {
	//sqls := "id|name|alias_name|scientific_name|shape|habit|harm|parasitic|distribution|govern_method|pest_image|"
	sqls := "id|title|date|paragraph|original|created|"
	sqls = strings.ReplaceAll(sqls, "|", ", ")
	fmt.Println(sqls)
}

func redisTest() {
	pestArray := make([]entry.PestTable, 0, 5)
	_ = config.InitRedisClient()
	err := config.GetJson("PestName-蝗虫", &pestArray)
	fmt.Println(err)
	fmt.Println(pestArray[0])

}

func v8Test() {
	table := make([]entry.PestSimilarity, 0)
	for i := 1; i < 100; i++ {
		var pt = make([]entry.PestTable, 0, 1)
		//pt := [1]entry.PestTable{}
		pt = append(pt, entry.PestTable{
			Name:      "13",
			AliasName: "24",
		})

		table = append(table, entry.PestSimilarity{
			PestName:   fmt.Sprintf("%d", i),
			Similarity: fmt.Sprintf("%d", i+i^i),
			Pests:      pt,
		})
	}
	pestTable := table[0].Pests
	fmt.Println(table[0].PestName)
	fmt.Println("lalalla", pestTable[0].Name)
}

func sile() {
	table := make([]int, 0)
	for i := 1; i < 100; i++ {
		table = append(table, i)
	}
	fmt.Println(table)
}

// 四选一
func isFourOne(choice string) bool {
	if len(choice) == 0 {
		return false
	}
	choice = strings.ToLower(choice)
	switch choice {
	case "pest", "weeds", "pesticide", "estimate":
		return true
	}
	return false
}

// 详细单页
func detailed() {
	url := "http://www.51agritech.com/pestweb/sc/SCHC/LYCLHC/WJZGYA.htm"
	d, _ := crawling.Detailed(url)
	fmt.Println(1, d.Name)
	fmt.Println(2, d.AliasName)
	fmt.Println(3, d.ScientificName)
	fmt.Println(4, d.Shape)
	fmt.Println(5, d.Habit)
	fmt.Println(6, d.Harm)
	fmt.Println(7, d.Parasitic)
	fmt.Println(8, d.Distribution)
	fmt.Println(9, d.GovernMethod)
	fmt.Println(10, d.PestImage)
}

// 保存图片
func saveImageTest() {
	url := "http://www.51agritech.com/pestweb/sc/SCHC/SSSCHC/DJM.jpg"
	save := `C:\Users\lwm\Desktop\1\`
	util.DownImage(save, url)
}

func t() []int {
	slice1 := make([]int, 10)
	fmt.Println(slice1)
	return slice1
}

func toBean() {

	str := `{"up load_part":[["蝼蛄","步甲","蝉","甘薯天蛾","食蚜蝇"],["77.09%","14.09%","1.06%","0.84%","0.60%"]]}`
	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.TypeOf(tempMap["up load_part"]).Kind())

	t := tempMap["up load_part"].(types.Slice)

	fmt.Println(t)
}

func getFileName() {
	srcFile, err := os.Open("C:/Users/lwm/Desktop/1/2.jpg")
	fmt.Println(err)
	fmt.Println(srcFile)
	dir, file := filepath.Split(srcFile.Name())
	fmt.Println(dir, file)
}

func post() {
	url := "http://121.36.19.108:8081/upload"

	// 具有Read和Write方法的可变大小的字节缓冲区
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	// 创建表单
	ContentType := writer.FormDataContentType()

	// 读取文件
	srcFile, err := os.Open("C:/Users/lwm/Desktop/1/2.jpg")
	if err != nil {
		fmt.Println("打开文件失败", err.Error())
	}
	defer srcFile.Close()
	_, fileName := filepath.Split(srcFile.Name())

	// 创建字段
	formFile, err := writer.CreateFormFile("file", fileName) // 提供表单中的字段名<img>和文件名<new.jpg>,返回值是可写的接口io.Writer
	if err != nil {
		fmt.Println("创建文件字段失败", err.Error())
		return
	}

	// 从文件读取数据，写入表单
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		fmt.Println("将文件拷贝到表单失败", err.Error())
	}

	// 发送
	writer.Close() // 发送之前必须调用Close()以写入结尾行
	resp, err := http.Post(url, ContentType, buf)
	if err != nil {
		fmt.Println("请求URL失败", url, err.Error())
		return
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}

	fmt.Println(string(b))
}
