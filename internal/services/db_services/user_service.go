package db_services

import (
	"errors"

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

func (us *UserService) Create(user *models.User) error {
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	passwordHashStr, err := utils.Password2hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = passwordHashStr
	if user.Status == 0 {
		user.Status = 1
	}
	return us.db.Create(user).Error
}

func (us *UserService) Update(user *models.User) error {
	var existing models.User
	if err := us.db.First(&existing, user.ID).Error; err != nil {
		return err
	}

	if user.Username != "" {
		existing.Username = user.Username
	}
	if user.Password != "" {
		passwordHashStr, err := utils.Password2hash(user.Password)
		if err != nil {
			return err
		}
		existing.Password = passwordHashStr
	}
	if user.Status != 0 {
		existing.Status = user.Status
	}

	if err := us.db.Save(&existing).Error; err != nil {
		return err
	}
	*user = existing
	return nil
}

func (us *UserService) GetUserByUsernamePassword(username, password string) (*models.User, error) {
	var info models.User
	passwordHashStr, password_err := utils.Password2hash(password)
	if password_err != nil {
		return nil, password_err
	}
	err := us.db.Where("username = ? and password = ? and status = 1", username, passwordHashStr).First(&info).Error
	return &info, err
}
