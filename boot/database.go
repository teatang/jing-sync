package boot

import (
	"jing-sync/models"
	"jing-sync/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	if !utils.FileExists("data/jing-sync.db") {
		DB, err = gorm.Open(sqlite.Open("data/jing-sync.db"), &gorm.Config{})
		DB.AutoMigrate(&models.User{}, &models.Engine{}, &models.Job{})
	} else {
		DB, err = gorm.Open(sqlite.Open("data/jing-sync.db"), &gorm.Config{})
	}
	if err != nil {
		panic("failed to connect database")
	}
}
