package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	// 示例使用
	// resp, err := Request("get", "https://httpbin.org/get", nil)
	// if err != nil {
	// 	fmt.Println("GET请求错误:", err)
	// } else {
	// 	fmt.Println("GET响应:", resp)
	// }

	data := map[string]interface{}{"requests": []interface{}{map[string]interface{}{
		"method": "GET",
		"url":    "https://httpbin.org/get",
	},"req"}}
	resp, err := Request("post", "https://www.act.mitsui-kinzoku.co.jp/wp-json/batch/v1", data)
	if err != nil {
		fmt.Printf("POST请求错误:%s\n", err)
	} else {
		fmt.Printf("POST响应: %s\n", resp)
	}
}

// Request 统一请求方法，支持get/post关键字选择
func Request(method string, url string, data interface{}) (string, error) {
	method = strings.ToLower(method)

	switch method {
	case "get":
		return getRequest(url)
	case "post":
		return postRequest(url, data)
	default:
		return "", fmt.Errorf("不支持的HTTP方法: %s", method)
	}
}

// getRequest 发送GET请求
func getRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// postRequest 发送POST请求
func postRequest(url string, data interface{}) (string, error) {
	// 序列化请求体
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
