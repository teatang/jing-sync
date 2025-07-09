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
		api.POST("/users", userController.CreateUser)
		api.GET("/users", userController.GetUsers)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)
	}

	return r
}
