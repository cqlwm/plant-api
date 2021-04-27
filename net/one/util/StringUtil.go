package util

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strings"
)

// 获取UUID
func Uuid(name string) string {
	v3 := uuid.NewV3(uuid.NewV1(), name)
	v3Str := v3.String()
	replace := strings.Replace(v3Str, "-", "", len(v3Str))
	return replace
}

// 获取文件扩展名
func Extension(fileName string) string {
	lastIndex := strings.LastIndex(fileName, ".")
	if lastIndex == -1 {
		return ""
	}
	k := fileName[lastIndex:]
	return k
}

func IsImage(name string) (img string, ok bool) {
	if len(name) == 0 {
		return "", false
	}
	name = Extension(name)
	if len(name) == 0 {
		return "", false
	}
	images := map[string]bool{
		".jpg": true,
		".png": true,
	}
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	_, ok = images[name]
	return name, ok
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
