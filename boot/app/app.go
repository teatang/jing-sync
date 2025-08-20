package app

import (
	"jing-sync/api/middlewares"
	"jing-sync/api/route"
	"jing-sync/boot/config"
	"jing-sync/boot/logger"

	"github.com/gin-gonic/gin"
	"time"
)

func WebSet() *gin.Engine {
	r := gin.Default()
	// 设置静态文件目录
	r.Static("/assets", "./frontend/dist/assets")

	// 设置timeout中间件
	config_timeout := config.Cfg.Timeout
	logger.GetLogger().Infof("web-site timeout: %d minutes", config_timeout)
	r.Use(middlewares.Timeout(time.Duration(config_timeout) * time.Second))
	// 设置日志中间件
	r.Use(middlewares.LoggerMiddleware())

	// 网站首页
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// 设置路由
	route.NewLoginRoute(r)
	route.NewProtectedRoute(r)

	return r
}
