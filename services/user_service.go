package services

import (
	"jing-sync/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) Create(user *models.User) error {
	return us.db.Create(user).Error
}

func (us *UserService) GetByID(id string) (*models.User, error) {
	var user models.User
	err := us.db.First(&user, id).Error
	return &user, err
}

func (us *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	err := us.db.Find(&users).Error
	return users, err
}

func (us *UserService) Update(id string, user *models.User) error {
	return us.db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

func (us *UserService) Delete(id string) error {
	return us.db.Delete(&models.User{}, id).Error
}
