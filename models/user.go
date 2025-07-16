package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    uint8     `json:"status" gorm:"default:1"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
}
