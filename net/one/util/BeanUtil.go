package util

import (
	"encoding/json"
	"log"
)

func BeanTo(src interface{}, target interface{}) bool {
	srcJ, err := json.Marshal(src)
	if err != nil {
		log.Println(err)
		return false
	}

	err = json.Unmarshal(srcJ, target)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
