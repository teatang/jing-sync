package db_services

import (
	"gorm.io/gorm"
	
	"jing-sync/internal/models"
)

type EngineService struct {
	BaseService[models.Engine]
}

func NewEngineService(db *gorm.DB) *EngineService {
	return &EngineService{BaseService[models.Engine]{db: db}}
}
