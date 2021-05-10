package crawling

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"plant-api/net/one/config"
	"strconv"
	"strings"
	"sync"
)

var pmcArr []provinceMappingCity
var codeMap map[string]*CityOne2

// 互斥锁
var mutex sync.Mutex

func ToDay(code string) {

}

// 未来一周天气情况
func FutureWeek(code int) (*WeatherResult, error) {
	url := buildUrlByCode(code)
	if len(url) == 0 {
		return nil, errors.New("找不到这个地区")
	}
	resp, err := http.Get(url)
	if err != nil {
		// TODO 请求错误
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// TODO 文档解析错误
		log.Println(err)
		return nil, err
	}

	weekInfo := doc.Find("#hourValues")

	var state, noticeImg string
	var hour, day, noRain int

	// 未来24小时气象
	weather2s := make([]weather2, 0)
	weekInfo.Find("#day0 .hour3").Each(func(i int, selection *goquery.Selection) {
		day := selection.Text()
		w := info24time(day)
		src, _ := selection.Find("img").Attr("src")
		src = deImgSrc(src)
		w.Img = src
		if w.Rain > 0 && hour == 0 {
			// 未来多少小时会降雨
			hour = (i + 1) * 3
		}
		weather2s = append(weather2s, w)
	})

	// 未来一周的天气
	weather1s := make([]weather1, 0)
	doc.Find("#day7 .weather .weatherWrap ").Each(func(i int, selection *goquery.Selection) {
		// 日期
		date := selection.Find(".date").Text()
		date = dateStr(date)
		// 天气
		desc := selection.Find(".desc").Text()
		desc = itemStr(desc, "转")
		// 天气图标
		descImg, _ := selection.Find("img").First().Attr("src")
		descImg = descImage(descImg)
		// 湿度
		selector := fmt.Sprintf("#day%d > div:nth-child(1) > div:nth-child(8)", i)
		text := weekInfo.Find(selector).Text()
		h := humidityStr(text)
		// 气温
		tmp := selection.Find(".tmp").Text()
		tmp = itemStr(tmp, "-")

		// 未来几天会有雨
		index := strings.Index(desc, "雨")
		if index != -1 && day == 0 {
			day = i + 1
			state = desc
			noticeImg = descImg
		} else {
			noRain++
		}

		weather1s = append(weather1s, weather1{
			Humidity:    h,
			Temperature: tmp,
			State:       desc,
			StateImage:  descImg,
			Time:        date,
		})
	})

	// 降雨提醒
	notice := notice(state, hour, day, noRain)
	if len(noticeImg) != 0 {
		noticeImg = weather1s[0].StateImage
	}
	rainNoticeEntry := rainNotice{notice, noticeImg}

	info := realInfo(code)

	wr := WeatherResult{weather1s, weather2s, &rainNoticeEntry, info}

	return &wr, nil
}

/*
天气分析：
气温高：未来X时间内会持续高温、请提前预备洒水措施
气温冷: 未来X时间内会持续低温、请提前预备护理措施
降雨、强降雨、连续降雨、连续强降雨
未来一段时间内不会出现降雨、
*/
type dayTmp struct {
	DateWeek string

	// 白天
	State1         string
	StateImg1      string
	windDirection1 string  // 风向
	windScale1     string  // 强度
	Temperature1   float64 // 温度

	// 夜间
	Temperature2   float64 // 温度
	State2         string
	StateImg2      string
	WindDirection2 string // 风向
	WindScale2     string // 强度
}

/*
天气、天气温度、提示信息
*/

// 天气提醒\温度\状态，对应首页
func DocNotice(code int) *rainNotice {

	url := buildUrlByCode(code)
	if len(url) == 0 {
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		// TODO 请求错误
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		// TODO 文档解析错误
		log.Println(err)
		return nil
	}

	weekTmpArr := make([]*dayTmp, 0)
	doc.Find("#day7 .weather").Each(func(i int, selection *goquery.Selection) {
		oneTmp := buildDayTmp(selection)
		weekTmpArr = append(weekTmpArr, oneTmp)
	})

	oneDayMoreInfoArr := make([]*oneDayMoreInfo, 0)
	doc.Find("#hourValues > div > .clearfix").Each(func(i int, selection *goquery.Selection) {
		selection.Find(".hour3").Each(func(j int, selection2 *goquery.Selection) {
			t := selection2.Text()
			imgSrc, _ := selection2.Find("img").Attr("src")
			dayMoreInfo := build24TimeOneDay(t)
			dayMoreInfo.Img = deImgSrc(imgSrc)
			oneDayMoreInfoArr = append(oneDayMoreInfoArr, dayMoreInfo)
		})
	})

	// preTime 未来小时降雨
	var preDay, preTime, noRain int
	var preStri, img string

	// 未来一周
	for k, v := range weekTmpArr {
		rain, b := indexRain(v.State1, v.State2)
		if preDay == 0 && len(rain) != 0 {
			preStri = rain
			preDay = k + 1
			if b {
				img = v.StateImg1
			} else {
				img = v.StateImg2
			}
			break
		} else {
			noRain++
		}
	}

	// 未来24小时
	for k, v := range oneDayMoreInfoArr[:8] {
		if v.Precipitation > 0 && preTime == 0 {
			preTime = (k + 1) * 3
			rain, b := indexRain(weekTmpArr[0].State1, weekTmpArr[0].State2)
			preStri = rain
			if b {
				img = weekTmpArr[0].StateImg1
			} else {
				img = weekTmpArr[0].StateImg2
			}
			break
		}
	}

	s := notice(preStri, preTime, preDay, noRain)
	imgSrc := deImgSrc(img)

	return &rainNotice{
		Msg: s,
		Img: imgSrc,
	}

}

// 是否是雨天
func indexRain(morning, afternoon string) (string, bool) {
	m := strings.Index(morning, "雨")
	a := strings.Index(afternoon, "雨")
	if m != -1 && a != -1 {
		return fmt.Sprintf("%s转%s", morning, afternoon), true
	}
	if m != -1 {
		return morning, true
	}
	if a != -1 {
		return afternoon, false
	}
	return "", false
}

func build24TimeOneDay(s string) *oneDayMoreInfo {
	// 23:00  -  24.9℃  2.4m/s  东南风  962.3hPa  72.6%  77.7%
	s = strings.TrimSpace(s)
	ss := strings.Split(s, "  ")
	o := oneDayMoreInfo{
		Time:          ss[0],
		Precipitation: dePrecipitation(ss[1]),
		Temperature:   deTemperature(ss[2]),
		WindSpeed:     deWindSpeed(ss[3]),
		WindDirection: trim(ss[4]),
		Pressure:      dePrecipitation(ss[5]),
		Humidity:      deHumidity(ss[5]),
	}
	return &o
}

// 解析降雨字符串
func dePrecipitation(de string) float64 {
	return deBase(de, "mm")
}

// 解析气温字符串
func deTemperature(de string) float64 {
	return deBase(de, "℃")
}

// 解析风速字符串
func deWindSpeed(de string) float64 {
	return deBase(de, "m/s")
}

// 解析气压字符串
func dePressure(de string) float64 {
	return deBase(de, "hPa")
}

// 解析湿度字符串
func deHumidity(de string) float64 {
	return deBase(de, "%")
}

// 解析气压字符串
func deImgSrc(de string) string {
	// http://image.nmc.cn/assets/img/w/40x40/3/1.png
	index := strings.LastIndex(de, "/")
	index = strings.LastIndex(de[:index], "/")
	de = strings.ReplaceAll(de[index+1:], "/", "-")
	return de
}

// 字符串去空格
func trim(s string) string {
	return strings.TrimSpace(s)
}

// 字符串转浮点数
func deBase(de string, tar string) float64 {
	de = strings.ReplaceAll(de, tar, "")
	de = strings.TrimSpace(de)
	r, _ := strconv.ParseFloat(de, 64)
	return r
}

type oneDayMoreInfo struct {
	Time          string  // 时间
	Img           string  // 图片
	Precipitation float64 // 降水
	Temperature   float64
	WindSpeed     float64
	WindDirection string  // 风向
	Pressure      float64 // 气压
	Humidity      float64 // 湿度
}

func buildDayTmp(selection *goquery.Selection) *dayTmp {
	date := selection.Find(".date").Text()
	date = strings.TrimSpace(date)

	state := selection.Find(".desc").Text()
	s1, s2 := stateStr(state)

	windd := selection.Find(".windd").Text()
	w1, w2 := stateStr(windd)

	winds := selection.Find(".winds").Text()
	ws1, ws2 := stateStr(winds)

	tmps := selection.Find(".tmp").Text()
	t1, t2 := tmp(tmps)

	var img [2]string
	selection.Find(".weathericon").Each(func(i int, sItme *goquery.Selection) {
		src, _ := sItme.Find("img").Attr("src")
		img[i] = src
	})

	d := dayTmp{
		DateWeek:       date,
		State1:         s1,
		StateImg1:      img[0],
		windDirection1: w1,
		windScale1:     ws1,
		Temperature1:   t1,
		Temperature2:   t2,
		State2:         s2,
		StateImg2:      img[1],
		WindDirection2: w2,
		WindScale2:     ws2,
	}

	return &d
}

func tmp(s string) (float64, float64) {
	t1, t2 := stateStr(s)
	tn1 := strings.ReplaceAll(t1, "℃", "")
	tn2 := strings.ReplaceAll(t2, "℃", "")

	var i1, i2 float64
	i1, _ = strconv.ParseFloat(tn1, 64)
	i2, _ = strconv.ParseFloat(tn2, 64)

	return i1, i2
}

func stateStr(s string) (string, string) {
	s = strings.TrimSpace(s)
	stateA := strings.Split(s, "  ")
	var state1, state2 string
	if len(stateA) == 2 {
		state1 = stateA[0]
		state2 = stateA[1]
	} else {
		state1 = ""
		state2 = stateA[0]
	}
	return state1, state2
}

type WeatherResult struct {
	W1     []weather1
	W2     []weather2
	Notice *rainNotice
	Info   map[string]interface{}
}

type rainNotice struct {
	Msg string
	Img string
}

func notice(state string, hour, day, noRain int) string {
	if hour != 0 {
		return fmt.Sprintf("未来%d小时可能将会出现%s", hour, state)
	}
	if day >= 5 {
		return fmt.Sprintf("未来一周可能将会出现%s", state)
	}
	if day != 0 {
		return fmt.Sprintf("未来%d天可能将会出现%s", day, state)
	}
	if noRain > 5 {
		return fmt.Sprintf("未来一周可能不会出现降雨")
	}
	if noRain != 0 {
		return fmt.Sprintf("未来%d天可能不会出现降雨", noRain)
	}
	return ""
}

// 通过code获取爬取网站的url
func buildUrlByCode(code int) string {
	if len(pmcArr) == 0 {
		mutex.Lock()
		if len(pmcArr) == 0 {
			if loadJson() != nil {
				fmt.Println("loadJson")
				return ""
			}
		}
		mutex.Unlock()
	}

	p := codeMap[strconv.Itoa(code)]
	if p == nil {
		fmt.Println("codeMap not noe")
		return ""
	}

	url2 := getCitySimpleUrl(p.ProvinceCode, strconv.Itoa(code))
	r := fmt.Sprintf("http://www.nmc.cn%s", url2)
	return r
}

func loadJson() error {
	//
	filePtr, err := os.Open(config.CityJsonPath)
	//filePtr, err := os.Open(`net/one/config/city.json`)

	if err != nil {
		log.Printf("Open file failed [Err:%s]", err.Error())
		return err
	}
	defer filePtr.Close()

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&pmcArr)
	if err != nil {
		log.Println("Decoder failed", err.Error())
		return err
	}

	pcMap := make(map[string]*CityOne2)

	for _, p := range pmcArr {
		for _, c := range p.City {
			pcMap[c.Code] = &CityOne2{
				ProvinceCode: p.Code,
				Province:     p.ProvinceName,
				City:         c.City,
			}
		}
	}

	codeMap = pcMap
	return nil
}

func getCitySimpleUrl(pCode string, code string) string {
	url := fmt.Sprintf("http://www.nmc.cn/rest/province/" + pCode)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	cityOnes := make([]CityOne, 0)
	_ = json.Unmarshal(body, &cityOnes)

	for _, v := range cityOnes {
		if v.Code == code {
			return v.Url
		}
	}

	return ""
}

type CityOne struct {
	Code     string
	Province string
	City     string
	Url      string
}

type CityOne2 struct {
	ProvinceCode string
	Province     string
	City         string
}

type provinceMappingCity struct {
	ProvinceName string
	Code         string
	City         []city
}
type city struct {
	City string
	Code string
}

type weather1 struct {
	Humidity    float64 // 湿度
	Temperature string  // 温度区间
	State       string  // 天气
	StateImage  string  // 天气图标
	Time        string  // 时间
}

type weather2 struct {
	Rain        float64 // 降水
	Temperature float64 // 温度
	Img         string
	Time        string // 时间
}

type temperature24Time struct {
	Humidity      float64 // 湿度
	Pressure      int     // 气压
	Rain1h        float64 // 降水
	Temperature   float64 // 温度
	Time          string  // 时间
	windDirection string  // 风向
	windScale     string  // 强度
}

func humidityStr(h string) float64 {
	h = strings.ReplaceAll(h, "%", " ")
	h = strings.TrimSpace(h)
	split := strings.Split(h, " ")
	var sum float64 = 0
	for _, v := range split {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			f, _ := strconv.ParseFloat(v, 64)
			sum += f
		}
	}
	sum = sum / float64(len(split))
	return sum
}

func dateStr(date string) string {
	date = strings.TrimSpace(date)
	date = strings.Split(date, " ")[0]
	date = strings.ReplaceAll(date, "/", "月")
	date = fmt.Sprintf("%s日", date)
	return date
}

func itemStr(s, item string) string {
	trimSpace := strings.TrimSpace(s)
	split := strings.Split(trimSpace, "  ")
	if len(split) == 1 {
		return split[0]
	}
	if split[0] == split[1] {
		return split[0]
	}
	return fmt.Sprintf("%s%s%s", split[0], item, split[1])
}

func descImage(img string) string {
	img = strings.ReplaceAll(img, "http://image.nmc.cn/assets/img/w/40x40/", "")
	img = strings.ReplaceAll(img, "/", "-")
	return img
}

func info24time(s string) weather2 {
	// 14:00  2.5mm  30.6℃  3.7m/s  东南风  964.9hPa  59.9%  70%
	s = strings.TrimSpace(s)
	split := strings.Split(s, "  ")

	rain := strings.TrimSpace(split[1])
	rain = strings.ReplaceAll(rain, "mm", "")
	rainNum, _ := strconv.ParseFloat(rain, 64)

	tem := strings.TrimSpace(split[2])
	tem = strings.ReplaceAll(tem, "℃", "")
	temNum, _ := strconv.ParseFloat(tem, 64)

	return weather2{
		Rain:        rainNum,
		Temperature: temNum,
		Time:        split[0],
	}
}

// ========================================
func RealInfo(code int) map[string]interface{} {
	return realInfo(code)
}

func realInfo(code int) map[string]interface{} {
	url := fmt.Sprintf("http://www.nmc.cn/rest/weather?stationid=%d", code)
	resp, err := http.Get(url)
	if err != nil {
		// TODO 异常请求
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// TODO IO异常
		return nil
	}

	res := map[string]interface{}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		// TODO 解析异常
		return nil
	}

	// 实时信息
	simpleReal := privateReal(res)
	// 地点
	place := privatePlace(res)

	return map[string]interface{}{
		"real":  simpleReal,
		"place": place,
	}

}

func privateReal(src map[string]interface{}) map[string]interface{} {
	publishTime := realPublishTime(src) // 发布时间
	weather := realWeather(src)
	temperature := weather["temperature"].(float64) // temperature温度
	day, night := temperatureOneDay(src)            // 高低温度
	rain := weather["rain"].(float64)               // 降雨量
	humidity := weather["humidity"].(float64)       // humidity湿度
	state := weather["info"].(string)               // 晴、雨
	stateImg := fmt.Sprintf(`4-%s.png`, weather["img"].(string))

	wind := realWind(src)
	windDirection := wind["direct"].(string) // 风向
	windScale := wind["power"].(string)      // 强度

	realMap := make(map[string]interface{})
	realMap["publishTime"] = publishTime
	realMap["rain"] = rain
	realMap["temperature"] = temperature
	realMap["day2night"] = fmt.Sprintf("%s℃-%s℃", day, night)
	realMap["humidity"] = humidity
	realMap["state"] = state
	realMap["stateImg"] = stateImg
	realMap["windDirection"] = windDirection
	realMap["windScale"] = windScale
	return realMap
}

func privatePlace(src map[string]interface{}) map[string]string {
	station := realStation(src)
	province := station["province"].(string)
	city := station["city"].(string)
	return map[string]string{
		"province": province,
		"city":     city,
	}
}

// temperatureOneDay
func temperatureOneDay(src map[string]interface{}) (string, string) {
	detail := predictDetail(src)
	days := detail[0].(map[string]interface{})
	day := days["day"].(map[string]interface{})["weather"].(map[string]interface{})["temperature"].(string)
	night := days["night"].(map[string]interface{})["weather"].(map[string]interface{})["temperature"].(string)
	return day, night
}

func predictDetail(src map[string]interface{}) []interface{} {
	predict := predict(src)
	detail := predict["detail"].([]interface{})
	return detail
}

func predict(src map[string]interface{}) map[string]interface{} {
	data := data(src)
	predict := data["predict"].(map[string]interface{})
	return predict
}

// 发布时间
func realPublishTime(src map[string]interface{}) string {
	station := real(src)["publish_time"].(string)
	return station
}

// realWind
func realWind(src map[string]interface{}) map[string]interface{} {
	station := real(src)["wind"].(map[string]interface{})
	return station
}

// realWeather
func realWeather(src map[string]interface{}) map[string]interface{} {
	station := real(src)["weather"].(map[string]interface{})
	return station
}

// realStation
func realStation(src map[string]interface{}) map[string]interface{} {
	station := real(src)["station"].(map[string]interface{})
	return station
}

func real(src map[string]interface{}) map[string]interface{} {
	data := data(src)
	real := data["real"].(map[string]interface{})
	return real
}

func data(src map[string]interface{}) map[string]interface{} {
	data := src["data"].(map[string]interface{})
	return data
}
