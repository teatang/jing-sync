package boot

import (
	"jing-sync/config"
	"jing-sync/controllers"
	"jing-sync/logger"
	"jing-sync/middlewares"
	"time"

	"github.com/gin-gonic/gin"
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

	db := GetDB()
	indexController := controllers.NewIndexController(db)
	userController := controllers.NewUserController(db)
	engineController := controllers.NewEngineController(db)
	jobController := controllers.NewJobController(db)
	openListController := controllers.NewOpenListController(db)

	// 登陆路由
	r.POST("/api/login", indexController.Login)

	// 受保护路由组
	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/user", userController.CreateUser)
		api.GET("/user", userController.GetPageUsers)
		api.GET("/user/:id", userController.GetUser)
		api.PUT("/user", userController.UpdateUser)
		api.DELETE("/user", userController.DeleteUser)

		api.POST("/engine", engineController.CreateEngine)
		api.GET("/engine", engineController.GetPageEngines)
		api.GET("/engine/:id", engineController.GetEngine)
		api.PUT("/engine", engineController.UpdateEngine)
		api.DELETE("/engine", engineController.DeleteEngine)

		api.POST("/job", jobController.CreateJob)
		api.GET("/job", jobController.GetPageJobs)
		api.GET("/job/:id", jobController.GetJob)
		api.PUT("/job", jobController.UpdateJob)
		api.DELETE("/job", jobController.DeleteJob)

		api.GET("/open-list", openListController.GetPageOpenList)
	}

	return r
}

func timeoutResponse(c *gin.Context) {
	c.JSON(504, gin.H{
		"code": 504,
		"msg":  "请求处理超时",
	})
}
