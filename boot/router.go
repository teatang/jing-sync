package boot

import (
	"github.com/gin-gonic/gin"
	"jing-sync/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 设置静态文件目录
	r.Static("/assets", "./frontend/dist/assets")

	// 网站首页
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	userController := controllers.NewUserController(DB)
	engineController := controllers.NewEngineController(DB)
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
	}

	return r
}
