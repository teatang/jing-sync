package app

import (
	"io/fs"
	"jing-sync/api/middlewares"
	"jing-sync/api/route"
	"jing-sync/boot/config"
	"jing-sync/boot/logger"
	"jing-sync/public"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WebSet() *gin.Engine {
	r := gin.Default()

	// webRootFS 就代表了整个 public/web 目录的内容。
	webRootFS, err := fs.Sub(public.WebFiles, "web")
	if err != nil {
		logger.GetLogger().Fatalf("Failed to create sub filesystem for 'web' directory: %v", err)
	}

	assetsFS, err := fs.Sub(webRootFS, "assets")
	if err != nil {
		logger.GetLogger().Fatalf("Failed to create sub filesystem for 'assets' directory: %v", err)
	}
	// 配置 Gin 路由来服务静态文件。
	r.StaticFS("/assets", http.FS(assetsFS))

	// 设置i18n中间件
	r.Use(middlewares.I18nMiddleware())
	// 设置timeout中间件
	config_timeout := config.Cfg.Timeout
	logger.GetLogger().Infof("web-site timeout: %d minutes", config_timeout)
	r.Use(middlewares.Timeout(time.Duration(config_timeout) * time.Second))
	// 设置日志中间件
	r.Use(middlewares.LoggerMiddleware())

	// 网站首页
	r.GET("/", func(c *gin.Context) {
		indexHTMLBytes, err := fs.ReadFile(webRootFS, "index.html")
		if err != nil {
			logger.GetLogger().Panicf("Failed to load index.html: %v", err)
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTMLBytes)
	})

	// 设置路由
	route.NewLoginRoute(r)
	route.NewProtectedRoute(r)

	return r
}
