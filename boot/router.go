package boot

import (
	"fmt"
	"jing-sync/controllers"
	"jing-sync/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 设置静态文件目录
	r.Static("/assets", "./frontend/dist/assets")

	// 设置日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))
	r.Use(middlewares.Logger())

	// 网站首页
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	userController := controllers.NewUserController(DB)
	engineController := controllers.NewEngineController(DB)
	jobController := controllers.NewJobController(DB)
	// 定义API路由
	api := r.Group("/api")
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
	}

	return r
}
