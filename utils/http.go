package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"wcj-go-common/core"
)

func ReadBody(r *http.Request) string {
	content, err := io.ReadAll(r.Body) //把	body 内容读入字符串 s
	if err != nil {
		log.Printf("%v", err)
	}
	return string(content)
}

func BindObject(r *http.Request, v interface{}) {
	content, err := io.ReadAll(r.Body) //把	body 内容读入字符串 s
	if err != nil {
		log.Printf("%v", err)
	}
	core.JsonToObject(&content, v)
}

func RespStringWithInt(w http.ResponseWriter, rInt int) {
	resp := strconv.Itoa(rInt)
	_, err := w.Write([]byte(resp))
	if err != nil {
		return
	}
}

func CloseReadCloser(r io.ReadCloser) {
	err := r.Close()
	if err != nil {
		return
	}
}

func Get(url string, any interface{}) {
	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("请求出错:%s", err)
	}
	defer CloseReadCloser(resp.Body)
	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应数据出错:%s", err)
	}
	core.JsonToObject(&body, any)
}

func HttpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP GET请求失败:%s", err)
		return nil
	}
	defer CloseReadCloser(resp.Body)
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应内容失败:%s", err)
		return nil
	}
	return body
}

func HttpGetContent(url string) string {
	//fmt.Println("httpGet url:" + url)
	return string(HttpGet(url))
}

func HttpDownload(downloadUrl string, targetFilepath string) bool {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return false
	}
	defer CloseReadCloser(resp.Body)
	out, err := os.Create(targetFilepath)
	if err != nil {
		log.Printf("创建目标文件失败:%v", err)
		return false
	}
	defer core.Close(out)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("下载异常:%v", err)
		return false
	}
	return true
}

func HttpPostJson(url string, reqBody string) string {
	method := "POST"
	payload := strings.NewReader(reqBody)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer core.Close(&res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}
