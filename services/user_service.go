package services

import (
	"gorm.io/gorm"
	"jing-sync/models"
)

type UserService struct {
	BaseService[models.User]
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{BaseService[models.User]{db: db}}
}
