package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jing-sync/models"
	"jing-sync/services"
	"net/http"
	"strconv"
)

type EngineController struct {
	BaseController
	engineService *services.EngineService
}

func NewEngineController(db *gorm.DB) *EngineController {
	return &EngineController{
		engineService: services.NewEngineService(db),
	}
}

// CreateUser 创建用户
func (ec *EngineController) CreateEngine(c *gin.Context) {
	var engine models.Engine
	if err := c.ShouldBindJSON(&engine); err != nil {
		ec.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ec.engineService.Create(&engine); err != nil {
		ec.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	ec.Success(c, engine)
}

// GetUser 获取单个用户
func (ec *EngineController) GetEngine(c *gin.Context) {
	id := c.Param("id")
	engine, err := ec.engineService.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ec.Error(c, http.StatusNotFound, "User not found")
		} else {
			ec.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ec.Success(c, engine)
}

// GetPageUsers 分页获取用户列表
func (ec *EngineController) GetPageEngines(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi("10")
	engines, err := ec.engineService.GetPageList(page, size)
	if err != nil {
		ec.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	ec.Success(c, engines)
}

// UpdateUser 更新用户
func (ec *EngineController) UpdateUser(c *gin.Context) {
	var engine models.Engine
	if err := c.ShouldBindJSON(&engine); err != nil {
		ec.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ec.engineService.Update(&engine); err != nil {
		if err == gorm.ErrRecordNotFound {
			ec.Error(c, http.StatusNotFound, "User not found")
		} else {
			ec.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	ec.Success(c, engine)
}

// DeleteUser 删除用户
func (ec *EngineController) DeleteEngine(c *gin.Context) {
	var engine models.Engine
	if err := c.ShouldBindJSON(&engine); err != nil {
		ec.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	id := strconv.Itoa(int(engine.ID))
	if err := ec.engineService.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			ec.Error(c, http.StatusNotFound, "User not found")
		} else {
			ec.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	ec.Success(c, engine)
}
