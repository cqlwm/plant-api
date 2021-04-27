package util

import (
	"github.com/huichen/sego"
	"plant-api/net/one/config"
	"sync"
)

// 载入词典
var segmenter sego.Segmenter

// 词库是否加载
var isLoad = false

// 互斥锁
var mutex sync.Mutex

// 线程安全的加载类
func callLoad() {
	if !isLoad {
		mutex.Lock()
		if !isLoad {
			segmenter.LoadDictionary(config.SegmenterLoadDictionary)
			isLoad = true
		}
		mutex.Unlock()
	}
}

// 线程安全的分词
func SafeParticiple(text string) []string {
	if !isLoad {
		callLoad()
	}
	textByte := []byte(text)
	// 分词
	mutex.Lock()
	segments := segmenter.Segment(textByte)
	mutex.Unlock()

	// 生成结果集
	result := make([]string, 0)
	for _, v := range segments {
		if !eliminate(v.Token().Text()) {
			result = append(result, v.Token().Text())
		}
	}

	return result
}

var eliminateMap = map[string]bool{
	"什么":   true,
	"原因":   true,
	"怎么":   true,
	"怎么办":  true,
	"是什么":  true,
	"怎么回事": true,
	"回事":   true,
}

// 消词
// return true 淘汰
func eliminate(original string) bool {
	if len(original) <= 3 {
		return true
	}
	s := eliminateMap[original]
	return s
}
