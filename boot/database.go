package boot

import (
	"jing-sync/config"
	"jing-sync/models"
	"jing-sync/utils"

	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	if DB == nil {
		InitDB()
	}
	return DB
}

func InitDB() string {
	var password string
	var err error

	utils.EnsureDir("data")

	dbPath := fmt.Sprintf("data/%s", config.Cfg.DbName)
	if !utils.FileExists(dbPath) {
		DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		password = AutoMigrate(DB)
	} else {
		DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
