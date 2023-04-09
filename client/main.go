package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// 建立一個gzip壓縮的JSON請求
	data := map[string]string{
		"message": "Hello, server!",
	}
	jsonData, _ := json.Marshal(data)
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(jsonData)
	writer.Close()

	// 建立一個HTTP請求
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/data", &buf)
	req.Header.Set("Content-Encoding", "gzip") // 告知api傳送內容編碼是gzip
	req.Header.Set("Accept-Encoding", "gzip")  // 允許響應是gzip

	// 發送HTTP請求並處理回應
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("header: ", resp.Header)

	// 讀取並解壓縮回應
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	// 解析回應的JSON資料
	var responseData map[string]string
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		panic(err)
	}

	// 輸出回應
	fmt.Println(responseData["message"])
	fmt.Println(responseData)
}
