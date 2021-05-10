package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"plant-api/net/one/entry"
	"strconv"
	"strings"
)

// 	url := "http://121.36.19.108:8081/upload"
//  "C:/Users/lwm/Desktop/1/2.jpg"
// 适用于昆虫识别的请求方法
func Post(url string, file string) (string, error) {

	// 创建具有Read和Write方法的可变大小的字节缓冲区
	buf := new(bytes.Buffer)

	// 创建表单
	writer := multipart.NewWriter(buf)
	ContentType := writer.FormDataContentType()

	// 读取文件
	srcFile, err := os.Open(file)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return "", errors.New("打开文件失败")
	}
	defer srcFile.Close()

	// 获取文件名称
	_, fileName := filepath.Split(srcFile.Name())

	// 创建字段
	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return "", errors.New("创建字段失败")
	}

	// 文件数据写入表单
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return "", errors.New("将文件拷贝到表单失败")
	}

	// 发送之前必须调用Close()以写入结尾行
	writer.Close()

	// Post请求
	resp, err := http.Post(url, ContentType, buf)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return "", errors.New("请求URL失败")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return "", errors.New("获取不到响应体")
	}

	return string(b), nil
}

// 花草请求POST
func PostFlower(image string) (*entry.WeedFlowerResult, error) {
	srcByte, err := ioutil.ReadFile(image)
	if err != nil {
		log.Println("PostFlower srcByte ", err)
		return nil, err
	}
	p1 := base64.StdEncoding.EncodeToString(srcByte)
	p2 := url.QueryEscape(p1)
	res := "img_base64=" + p2

	c := &http.Client{}
	url := "http://plantgw.nongbangzhu.cn/plant/recognize2"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(res)))
	if err != nil {
		fmt.Println("PostFlower NewRequest ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", "APPCODE 942f9103297e4f508570e3dd2272bf40")
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	bodys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("响应：", err)
		return nil, err
	}
	flower := entry.WeedFlowerResult{}
	err = json.Unmarshal(bodys, &flower)
	if err != nil {
		return nil, err
	}

	fmt.Println("PostFlower OK ", flower)
	return &flower, nil
}

// 获取花草信息
func PostWeedsInfo() {
	info := "I2XadqHedSdQcZni"
	url := "http://plantgw.nongbangzhu.cn/plant/info"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader("code="+info))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", "APPCODE 942f9103297e4f508570e3dd2272bf40")
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("响应 ", err)
		return
	}
	fmt.Println(string(body))
}

// 病害识别
// http://senseagro.market.alicloudapi.com/api/senseApi
func Disease(crop_id string, image_url string) (*disResult, error) {
	ret, err := Disease0(crop_id, image_url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := buildResult(ret)
	return &result, nil
}
func Disease0(crop_id string, image_url string) (string, error) {
	url := "http://senseagro.market.alicloudapi.com/api/senseApi"

	client := &http.Client{}

	requestBody := fmt.Sprintf("crop_id=%s&image_url=%s", crop_id, image_url)
	log.Println("param ", requestBody)

	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Authorization", "APPCODE 942f9103297e4f508570e3dd2272bf40")
	if err != nil {
		log.Println("Disease(crop_id string, image_url string)  request err ", err)
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("get do ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("响应 ", err)
		return "", err
	}
	sBody := fmt.Sprintf("%s", string(body))
	//log.Println("request ok ", string(body))
	//result := disResult{}
	//_ = json.Unmarshal([]byte(sBody), &result)
	//if err != nil {
	//	log.Println("json fail", err.Error())
	//	return nil, err
	//}
	//result := buildResult(sBody)
	//log.Println(result)
	return sBody, nil
}

func buildResult(s string) disResult {
	// s := `{"status":"1","msg":"操作成功","content":{"result":"蚧壳虫","score":90.59}}`
	result := disResult{}
	_ = json.Unmarshal([]byte(s), &result)
	log.Println(result)
	return result
}

// {"status":"1","msg":"操作成功","content":{"result":"柑橘树脂病","score":94.46}}
type disResult struct {
	Status  string
	Msg     string
	Content disContent
}
type disContent struct {
	Result string
	Score  float64
}

// 为模型增加训练集

// 模型训练

// 以图搜索
func SearchByImage(filePath string) ([]int, error) {
	url := `http://47.96.7.148:35000/api/v1/search`

	fields := make(map[string]string)
	fields["Num"] = "75"
	var body, err = form("post", url, fields, "file", filePath)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resArr := [][]interface{}{}
	_ = json.Unmarshal(body, &resArr)

	iarr := make([]int, 0)
	for _, v := range resArr {
		imagename := v[0].(string)
		index1 := strings.LastIndex(imagename, "/")
		index2 := strings.LastIndex(imagename, ".")
		//one, _ := strconv.ParseInt(, 16, 64)
		one, _ := strconv.Atoi(imagename[index1+1 : index2])
		iarr = append(iarr, one)
	}

	return iarr, nil
}

func form(method string, url string, fields map[string]string,
	fileKey string, filePath string) ([]byte, error) {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	if len(fileKey) != 0 {
		file, err := os.Open(filePath)
		defer file.Close()
		part2, err := writer.CreateFormFile(fileKey, filepath.Base(filePath))
		_, err = io.Copy(part2, file)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil

	//resArr := [][]interface{}{}
	//_ = json.Unmarshal(body, &resArr)
	//
	//iarr := make([]int64, 0)
	//for _, v := range resArr {
	//	imagename := v[0].(string)
	//	index1 := strings.LastIndex(imagename, "/")
	//	index2 := strings.LastIndex(imagename, ".")
	//	one, _ := strconv.ParseInt(imagename[index1+1:index2], 16, 64)
	//	iarr = append(iarr, one)
	//}
	//
	//return iarr, nil
}
