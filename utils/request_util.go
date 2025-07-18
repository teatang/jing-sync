package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type RequestConfig struct {
	URL     string
	Method  string // GET/POST/PUT/DELETE
	Headers map[string]string
	Body    interface{} // 支持string/[]byte/struct
	Timeout time.Duration
}

type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func SendRequest(ctx context.Context, config *RequestConfig) (*Response, error) {
	// 处理请求体
	var body io.Reader
	switch v := config.Body.(type) {
	case nil:
		body = nil
	case string:
		body = bytes.NewBufferString(v)
	case []byte:
		body = bytes.NewBuffer(v)
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}
	if _, ok := config.Headers["Content-Type"]; !ok && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 配置客户端
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Header:     resp.Header,
	}, nil
}
