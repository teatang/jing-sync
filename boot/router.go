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
		api.GET("/user", userController.GetUsers)
		api.GET("/user/:id", userController.GetUser)
		api.PUT("/user/:id", userController.UpdateUser)
		api.DELETE("/user/:id", userController.DeleteUser)
	}

	return r
}
