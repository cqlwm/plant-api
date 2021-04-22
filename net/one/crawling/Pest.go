package crawling

import (
	"encoding/json"
	"errors"
	"fmt"
	"plant-api/net/one/util"
)

const (
	url = "http://121.36.19.108:8081/upload"
)

func PestPost(image string) ([][]string, error) {
	// 请求接口
	result, err := util.Post(url, image)
	if err != nil {
		return nil, err
	}

	var tempMap map[string]interface{}
	err = json.Unmarshal([]byte(result), &tempMap)
	if err != nil {
		fmt.Errorf(err.Error())
		return nil, errors.New("响应结果JSON反序列化失败")
	}

	var t interface{}
	var ok bool
	t, ok = tempMap["up load_part"]
	if !ok {
		t, ok = tempMap["upload_part"]
	}

	if !ok {
		// Unicode编码
		return nil, errors.New("解析结果失败")
	}

	// 获取最大可能
	rb, err := json.Marshal(t)
	if err != nil {
		// 序列化失败
		return nil, err
	}

	var resArray [][]string
	err = json.Unmarshal(rb, &resArray)
	if err != nil {
		// 反序列化失败
		return nil, err
	}

	return resArray, nil
}
