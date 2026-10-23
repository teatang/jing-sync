package controllers

import (
	"jing-sync/internal/models"
	"jing-sync/internal/services"
	"jing-sync/internal/services/db_services"
	"jing-sync/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EngineController struct {
	BaseController
	engineService *db_services.EngineService
	db            *gorm.DB
}

func NewEngineController(db *gorm.DB) *EngineController {
	return &EngineController{
		engineService: db_services.NewEngineService(db),
		db:            db,
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

type TestEngineRequest struct {
	ID    uint   `json:"id"`
	Url   string `json:"url"`
	Token string `json:"token"`
}

// TestEngine 测试引擎连通性
func (ec *EngineController) TestEngine(c *gin.Context) {
	var req TestEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ec.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var engine *models.Engine
	if req.ID > 0 {
		var err error
		engine, err = ec.engineService.GetByID(strconv.Itoa(int(req.ID)))
		if err != nil {
			ec.Error(c, http.StatusNotFound, "Engine not found")
			return
		}
	} else {
		engine = &models.Engine{
			Url:   req.Url,
			Token: req.Token,
		}
	}

	client := services.NewOpenListClientFromEngine(engine, ec.db)
	if err := client.TestConnection(); err != nil {
		ec.Error(c, http.StatusBadRequest, "Connection test failed: "+err.Error())
		return
	}

	ec.Success(c, gin.H{"status": "ok", "message": "Connection test passed"})
}
