package test

import (
	"jing-sync/internal/utils"

	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	transport := utils.NewHandlerTransport(handler)

	// 测试GET请求
	resp, err := utils.Request("GET", "http://mock.local", nil, &utils.RequestOption{
		Transport: transport,
	})
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	var result map[string]string
	json.Unmarshal(resp, &result)
	if result["status"] != "ok" {
		t.Errorf("响应数据异常: %s", string(resp))
	}
}

func TestPostWithRetry(t *testing.T) {
	retryCount := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if retryCount < 2 {
			retryCount++
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
	transport := utils.NewHandlerTransport(handler)

	// 测试重试机制
	_, err := utils.Request("POST", "http://mock.local", nil, &utils.RequestOption{
		Retry:     2,
		Timeout:   time.Second,
		Transport: transport,
	})
	if err != nil || retryCount != 2 {
		t.Errorf("重试机制异常 实际重试%d次", retryCount)
	}
}

func TestTimeout(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
	})
	transport := utils.NewHandlerTransport(handler)

	// 测试超时控制
	_, err := utils.Request("GET", "http://mock.local", nil, &utils.RequestOption{
		Timeout:   10 * time.Millisecond,
		Transport: transport,
	})
	if err == nil {
		t.Error("未触发超时限制")
	}
}
