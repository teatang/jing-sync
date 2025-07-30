package boot

import (
	"jing-sync/config"
	"jing-sync/models"
	"jing-sync/utils"

	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}

func InitDB() string {
	var password string
	var err error

	utils.EnsureDir("data")

	dbPath := fmt.Sprintf("data/%s", config.Cfg.DbName)
	if !utils.FileExists(dbPath) {
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		password = AutoMigrate(db)
	} else {
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}

	return password
}

func AutoMigrate(db *gorm.DB) string {
	db.AutoMigrate(&models.User{}, &models.Engine{}, &models.Job{})
	password, err := utils.SecureRandString(10)
	if err != nil {
		panic(err)
	}
	passwordHashStr, err := utils.Password2hash(password)
	if err != nil {
		panic(err)
	}
	db.Create(&models.User{Username: "admin", Password: passwordHashStr})

	return password
}
