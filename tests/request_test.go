package tests

import (
	"encoding/json"
	"jing-sync/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRequest(t *testing.T) {
	// 创建模拟服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	// 测试GET请求
	resp, err := utils.Request("GET", ts.URL, nil, nil)
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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if retryCount < 2 {
			retryCount++
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	// 测试重试机制
	_, err := utils.Request("POST", ts.URL, nil, &utils.RequestOption{
		Retry:   3,
		Timeout: time.Second,
	})
	if err != nil || retryCount != 2 {
		t.Errorf("重试机制异常 实际重试%d次", retryCount)
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 故意延迟
	}))
	defer ts.Close()

	// 测试超时控制
	_, err := utils.Request("GET", ts.URL, nil, &utils.RequestOption{
		Timeout: time.Second,
	})
	if err == nil {
		t.Error("未触发超时限制")
	}
}
