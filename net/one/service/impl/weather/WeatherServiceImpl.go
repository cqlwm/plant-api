package weather

import (
	"plant-api/net/one/crawling"
)

type WeatherServiceImpl struct{}

func (ws *WeatherServiceImpl) Info(code int) (*crawling.WeatherResult, error) {
	week, err := crawling.FutureWeek(code)
	return week, err
}

func (ws *WeatherServiceImpl) Simple(code int) map[string]interface{} {
	rainNotice := crawling.DocNotice(code)
	info := crawling.RealInfo(code)

	ret := make(map[string]interface{}, 0)
	ret["rainNotice"] = rainNotice
	ret["temperature"] = info["real"].(map[string]interface{})["temperature"]
	ret["stateImg"] = info["real"].(map[string]interface{})["stateImg"]
	ret["state"] = info["real"].(map[string]interface{})["state"]
	return ret
}
