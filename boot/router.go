package boot

import (
	"jing-sync/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	userController := controllers.NewUserController(DB)

	api := r.Group("/api")
	{
		api.POST("/user", userController.CreateUser)
		api.GET("/user", userController.GetPageUsers)
		api.GET("/user/:id", userController.GetUser)
		api.PUT("/user", userController.UpdateUser)
		api.DELETE("/user", userController.DeleteUser)
	}

	return r
}
