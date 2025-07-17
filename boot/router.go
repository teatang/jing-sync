package boot

import (
	"github.com/gin-gonic/gin"
	"jing-sync/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	userController := controllers.NewUserController(DB)
	engineController := controllers.NewEngineController(DB)

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
		api.PUT("/engine", engineController.UpdateUser)
		api.DELETE("/engine", engineController.DeleteEngine)
	}

	return r
}
