package route

import (
	"jing-sync/api/controllers"
	"jing-sync/boot/database"

	"github.com/gin-gonic/gin"
)

func NewLoginRoute(r *gin.Engine) {
	db := database.GetDB()
	indexController := controllers.NewIndexController(db)

	// 登陆路由
	r.POST("/api/login", indexController.Login)
}
