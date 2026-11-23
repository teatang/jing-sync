package utils

import (
	"jing-sync/boot/logger"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

type RequestOption struct {
	Headers   map[string]string
	Timeout   time.Duration // 默认30秒
	Retry     int           // 默认不重试
	Transport http.RoundTripper
}

var customTransport http.RoundTripper

// SetCustomTransport 设置自定义 RoundTripper，供单元测试等场景 mock 网络请求
func SetCustomTransport(t http.RoundTripper) {
	customTransport = t
}

// HandlerTransport 将 http.Handler 直接转为 RoundTripper，避免真实网络调用
type HandlerTransport struct {
	Handler http.Handler
}

func (ht *HandlerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	done := make(chan struct{})
	rec := httptest.NewRecorder()
	go func() {
		ht.Handler.ServeHTTP(rec, req)
		close(done)
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case <-done:
		return rec.Result(), nil
	}
}

// NewHandlerTransport 创建基于 Handler 的内存传输层
func NewHandlerTransport(h http.Handler) *HandlerTransport {
	return &HandlerTransport{Handler: h}
}

func Request(method, url string, body interface{}, opt *RequestOption) ([]byte, error) {
	if opt == nil {
		opt = &RequestOption{Timeout: 30 * time.Second}
	}

	var reader io.Reader
	if body != nil {
		switch v := body.(type) {
		case []byte:
			reader = bytes.NewBuffer(v)
		case string:
			reader = bytes.NewBufferString(v)
		default:
			data, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewBuffer(data)
		}
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}
	if _, ok := opt.Headers["Content-Type"]; !ok && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: opt.Timeout}
	if opt.Transport != nil {
		client.Transport = opt.Transport
	} else if customTransport != nil {
		client.Transport = customTransport
	}
	var resp *http.Response

	// 重试逻辑
	for i := 0; i <= opt.Retry; i++ {
		resp, err = client.Do(req)
		//记录日志
		logger.GetLogger().Infof("Request: method=%s, url=%s, body=%v, opt=%+v, resp=%+v, err=%+v", method, url, body, opt, resp, err)
		if err == nil && resp.StatusCode < 500 {
			break
		}
		if i < opt.Retry {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func Get(url string, opt *RequestOption) ([]byte, error) {
	return Request(http.MethodGet, url, nil, opt)
}

func Post(url string, body interface{}, opt *RequestOption) ([]byte, error) {
	return Request(http.MethodPost, url, body, opt)
}
