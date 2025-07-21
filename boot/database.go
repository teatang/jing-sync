package boot

import (
	"jing-sync/models"
	"jing-sync/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() string {
	var password string
	var err error

	utils.EnsureDir("data")

	if !utils.FileExists("data/jing-sync.db") {
		DB, err = gorm.Open(sqlite.Open("data/jing-sync.db"), &gorm.Config{})
		password = AutoMigrate(DB)
	} else {
		DB, err = gorm.Open(sqlite.Open("data/jing-sync.db"), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}

	return password
}

func AutoMigrate(DB *gorm.DB) string {
	DB.AutoMigrate(&models.User{}, &models.Engine{}, &models.Job{})
	password, err := utils.SecureRandString(10)
	if err != nil {
		panic(err)
	}
	passwordHashStr, err := utils.Password2hash(password)
	if err != nil {
		panic(err)
	}
	DB.Create(&models.User{Username: "admin", Password: passwordHashStr})

	return password
}
