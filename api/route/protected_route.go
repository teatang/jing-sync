package route

import (
	"jing-sync/api/controllers"
	"jing-sync/api/middlewares"
	"jing-sync/boot/database"

	"github.com/gin-gonic/gin"
)

func NewProtectedRoute(r *gin.Engine) {
	db := database.GetDB()
	userController := controllers.NewUserController(db)
	engineController := controllers.NewEngineController(db)
	jobController := controllers.NewJobController(db)
	openListController := controllers.NewOpenListController(db)
	syncLogController := controllers.NewSyncLogController(db)

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

		// 同步日志相关
		api.GET("/job/:job_id/sync-logs", syncLogController.GetSyncLogsByJobId)
		api.GET("/sync-log/:id", syncLogController.GetSyncLog)
		api.POST("/job/:job_id/sync", syncLogController.TriggerSyncManually)
	}
}
