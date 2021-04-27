package util

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// save 保存的磁盘路径
// url 网络路径
func DownImage(save string, url string) {
	// 请求图片
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()

	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	fileName := path.Base(url)
	file, err := os.Create(save + fileName)
	if err != nil {
		fmt.Println("创建图片失败")
		return
	}
	defer file.Close()

	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	//保存文件
	_, _ = io.Copy(writer, reader)
	//fmt.Printf("Total length: %d", written)
}

//
func DownImage2(saveImage string, urlImage string) {
	// 请求图片
	res, err := http.Get(urlImage)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()

	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(saveImage)
	if err != nil {
		fmt.Println("创建图片失败")
		return
	}
	defer file.Close()

	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	//保存文件
	_, _ = io.Copy(writer, reader)
	//fmt.Printf("Total length: %d", written)
}
