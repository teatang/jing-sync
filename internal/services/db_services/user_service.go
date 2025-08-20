package db_services

import (
	"gorm.io/gorm"

	"jing-sync/internal/models"
	"jing-sync/internal/utils"
)

type UserService struct {
	BaseService[models.User]
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{BaseService[models.User]{db: db}}
}

func (UserService *UserService) GetUserByUsernamePassword(username, password string) (*models.User, error) {
	var info models.User
	passwordHashStr, password_err := utils.Password2hash(password)
	if password_err != nil {
		return nil, password_err
	}
	err := UserService.db.Where("username = ? and password = ?", username, passwordHashStr).First(&info).Error
	return &info, err
}
