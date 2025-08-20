package controllers

import (
	"jing-sync/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type OpenListController struct {
	BaseController
	openListService *services.OpenListService
}

func NewOpenListController(db *gorm.DB) *OpenListController {
	return &OpenListController{openListService: services.NewOpenListService(db)}
}

// GetPageJobs 分页获取用户列表
func (olc *OpenListController) GetPageOpenList(c *gin.Context) {
	engine_id := c.Query("engine_id")
	path := c.DefaultQuery("path", "/")
	infos, err := olc.openListService.GetOpenListInfo(engine_id, path)
	if err != nil {
		olc.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	olc.Success(c, infos)
}
