package db_services

import (
	"gorm.io/gorm"
	
	"jing-sync/internal/models"
)

type JobService struct {
	BaseService[models.Job]
}

func NewJobService(db *gorm.DB) *JobService {
	return &JobService{BaseService[models.Job]{db: db}}
}
