package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jing-sync/models"
	"jing-sync/utils"
	"jing-sync/services/db_services"
	"net/http"
	"strconv"
)

type JobController struct {
	BaseController
	jobService *db_services.JobService
}

func NewJobController(db *gorm.DB) *JobController {
	return &JobController{
		jobService: db_services.NewJobService(db),
	}
}

// CreateJob 创建用户
func (uc *JobController) CreateJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		uc.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	job.UserId = utils.GetTokenUserId(c)

	if err := uc.jobService.Create(&job); err != nil {
		uc.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	uc.Success(c, job)
}

// GetJob 获取单个用户
func (uc *JobController) GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := uc.jobService.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Error(c, http.StatusNotFound, "Job not found")
		} else {
			uc.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	uc.Success(c, job)
}

// GetPageJobs 分页获取用户列表
func (uc *JobController) GetPageJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi("10")
	jobs, err := uc.jobService.GetPageList(page, size)
	if err != nil {
		uc.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	uc.Success(c, jobs)
}

// UpdateJob 更新用户
func (uc *JobController) UpdateJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		uc.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := uc.jobService.Update(&job); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Error(c, http.StatusNotFound, "Job not found")
		} else {
			uc.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	uc.Success(c, job)
}

// DeleteJob 删除用户
func (uc *JobController) DeleteJob(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		uc.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	id := strconv.Itoa(int(job.ID))
	if err := uc.jobService.Delete(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Error(c, http.StatusNotFound, "Job not found")
		} else {
			uc.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	uc.Success(c, job)
}
