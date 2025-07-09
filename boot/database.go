package boot

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jing-sync/utils"
	"jing-sync/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	if !utils.FileExists("data/user.db") {
	    DB, err = gorm.Open(sqlite.Open("data/user.db"), &gorm.Config{})
		DB.AutoMigrate(&models.User{})
	} else {
		DB, err = gorm.Open(sqlite.Open("data/user.db"), &gorm.Config{})
	}
	if err != nil {
		panic("failed to connect database")
	}
}
