package middlewares

import (
	"jing-sync/logger"
	"time"
	"bytes"
	"io"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 读取原始请求的 Body 数据
        bodyBytes, _ := io.ReadAll(c.Request.Body)
        // 将原始的 Body 数据重新放入请求的 Body 中，以便后续的处理函数可以再次读取（重置）
        c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		// 将读取到的 Body 数据转换为 map 类型
		var bodyMap map[string]interface{}
		json.Unmarshal(bodyBytes, &bodyMap)

		c.Next()

		logger.Log.WithFields(logrus.Fields{
			"status":       c.Writer.Status(),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"get_params":   c.Request.URL.Query(),
			"request_body": bodyMap,
			"client_ip":    c.ClientIP(),
			"cost_time":    time.Since(start).Milliseconds(),
			"user_agent":   c.Request.UserAgent(),
		}).Info("http_request")
	}
}
