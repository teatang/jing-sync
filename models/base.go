package models

import "time"

type BaseModel struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
	Status     uint8     `json:"status" gorm:"default:1"`
}
