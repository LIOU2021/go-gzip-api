package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/data", handleRequest)
	host := "127.0.0.1:8080"
	fmt.Printf("run: %s\n", host)
	http.ListenAndServe(host, nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 檢查請求是否為gzip格式
	if r.Header.Get("Content-Encoding") == "gzip" {
		// 解壓縮gzip
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
			return
		}
		defer reader.Close()
		r.Body = io.NopCloser(reader)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to parse request", http.StatusInternalServerError)
		return
	}

	// 處理請求並生成回應
	data := map[string]string{
		"message": "Hello, world!",
		"request": string(b),
	}
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	// 壓縮JSON回應
	w.Header().Set("Content-Encoding", "gzip")
	writer := gzip.NewWriter(w)
	defer writer.Close()
	// 壓縮gzip
	writer.Write(response)
}
