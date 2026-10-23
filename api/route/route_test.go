package route

import (
	"testing"

	"jing-sync/boot/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRouteRegistrationNoConflict 验证路由注册不会因为 Gin 路由树通配符冲突而 panic
func TestRouteRegistrationNoConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.SetDefaultConfig(config.EnvTypeTest)
	config.Cfg.DbName = ":memory:"
	r := gin.New()

	// 验证路由注册正常无 panic
	assert.NotPanics(t, func() {
		NewLoginRoute(r)
		NewProtectedRoute(r)
	})
}

