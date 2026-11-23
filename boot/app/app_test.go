package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"jing-sync/boot/config"
	"jing-sync/boot/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSetInitialization(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("data")
	})
	gin.SetMode(gin.TestMode)
	logger.LoggerInit()
	config.SetDefaultConfig(config.EnvTypeTest)
	config.Cfg.DbName = ":memory:"

	engine := WebSet()
	require.NotNil(t, engine)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/login", nil)
	engine.ServeHTTP(w, req)
	// 验证路由已匹配（即便没有 payload 也不会返回 404）
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}
