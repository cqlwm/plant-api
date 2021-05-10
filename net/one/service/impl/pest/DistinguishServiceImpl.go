package pest

import (
	"log"
	"plant-api/net/one/config"
	"plant-api/net/one/crawling"
	"plant-api/net/one/entry"
	"plant-api/net/one/util"
	"strings"
)

// 识别
var typeMap = map[string]string{
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

type DistinguishServiceImpl struct{}

// 识别害虫服务
func (dsi *DistinguishServiceImpl) Pest(pestImage string) ([]entry.PestSimilarity, error) {
	// 识别图片[["蝗虫"]["98%"]]
	res, err := crawling.PestPost(pestImage)
	if err != nil {
		return nil, err
	}

	var pestSimilarityArray = make([]entry.PestSimilarity, 0)

	// 为每种害虫匹配防止方案
	for i := 0; i < len(res[0]); i++ {
		pestArray := FindPestByName(res[0][i])
		p := entry.PestSimilarity{
			PestName:   res[0][i],
			Similarity: res[1][i],
			Pests:      pestArray,
		}
		pestSimilarityArray = append(pestSimilarityArray, p)
	}

	// TODO 根据用户，生成查询日志
	return pestSimilarityArray, nil
}

// 杂草识别服务
func (dsi *DistinguishServiceImpl) Weeds(pestImage string) {
	//
}

// 病害识别服务
func (dsi *DistinguishServiceImpl) Disease(pestImage string, completeNewName string) (*entry.WeedResult, error) {
	res, err := util.PostFlower(pestImage)
	log.Println("PostFlower  ", res)
	// 辨认果木
	key := ""
	if err == nil {
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
	}

	if key == "" {
		key = "1"
	}

	// http://graceful.top/plant/image/item/GJ01.jpg
	//testUrl := `http://graceful.top/plant/image/item/GJ01.jpg`
	//dire, err := util.Disease(key, testUrl)

	dire, err := util.Disease(key, config.ReUrlImage+completeNewName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	diseaseName := dire.Content.Result

	log.Println("Disease diseaseName ", diseaseName)

	// 百度百科爬虫
	//config.MongoDBConn()
	//config.DeleteAll(config.WeedColl)
	find := crawling.Find(diseaseName)
	return find, nil
}
