package controllers

import (
	"jing-sync/internal/models"
	"jing-sync/internal/utils"
	"jing-sync/internal/services/db_services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EngineController struct {
	BaseController
	engineService *db_services.EngineService
}

func NewEngineController(db *gorm.DB) *EngineController {
	return &EngineController{
		engineService: db_services.NewEngineService(db),
	}
}

func (ec *EngineController) CreateEngine(c *gin.Context) {
	var engine models.Engine
	if err := c.ShouldBindJSON(&engine); err != nil {
		ec.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	engine.UserId = utils.GetTokenUserId(c)

	if err := ec.engineService.Create(&engine); err != nil {
		ec.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	ec.Success(c, engine)
}

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

func (ec *EngineController) UpdateEngine(c *gin.Context) {
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
