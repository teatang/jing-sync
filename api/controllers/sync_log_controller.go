package controllers

import (
	"jing-sync/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SyncLogController struct {
	BaseController
	syncService *services.SyncService
}

func NewSyncLogController(db *gorm.DB) *SyncLogController {
	return &SyncLogController{
		syncService: services.NewSyncService(db),
	}
}

// GetSyncLogsByJobId 获取任务的同步日志列表
func (c *SyncLogController) GetSyncLogsByJobId(ctx *gin.Context) {
	jobId, err := strconv.ParseUint(ctx.Param("job_id"), 10, 32)
	if err != nil {
		c.Error(ctx, http.StatusBadRequest, "Invalid job ID")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	logs, err := c.syncService.GetSyncLogsByJobId(uint(jobId), page, size)
	if err != nil {
		c.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	c.Success(ctx, logs)
}

// GetSyncLog 获取单条同步日志
func (c *SyncLogController) GetSyncLog(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		c.Error(ctx, http.StatusBadRequest, "Invalid log ID")
		return
	}

	log, err := c.syncService.GetSyncLogById(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(ctx, http.StatusNotFound, "Sync log not found")
		} else {
			c.Error(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Success(ctx, log)
}

// TriggerSyncManually 手动触发同步
func (c *SyncLogController) TriggerSyncManually(ctx *gin.Context) {
	jobId, err := strconv.ParseUint(ctx.Param("job_id"), 10, 32)
	if err != nil {
		c.Error(ctx, http.StatusBadRequest, "Invalid job ID")
		return
	}

	syncLog, err := c.syncService.ExecuteSync(uint(jobId))
	if err != nil {
		c.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	c.Success(ctx, syncLog)
}
