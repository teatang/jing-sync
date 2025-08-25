package middlewares

import (
	bootI18n "jing-sync/boot/i18n"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// JWT认证中间件
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		i18nBundle := bootI18n.GetI18nBundle()
		lang := c.GetHeader("Accept-Language")
		localizer := i18n.NewLocalizer(i18nBundle, lang)
		c.Set("localizer", localizer)
		c.Next()
	}
}
