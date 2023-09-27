package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ChatHost          = "http://aichat.tanyadong.com:1234"
	FantasticChatHost = "http://fanchat.tanyaodng.com:8080"
)

func DoPost(url string, json string) (string, error) {
	payload := strings.NewReader(json)

	// 发送POST请求
	fmt.Println("url = " + url)
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		fmt.Println("POST请求发送失败:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应数据失败:", err)
		return "", err
	}

	fmt.Println("响应状态码:", resp.Status)
	fmt.Println("响应内容:", string(body))
	return string(body), nil
}
