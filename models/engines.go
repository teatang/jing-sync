package models

import "time"

type Engine struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    uint8     `json:"status" gorm:"default:1"`
	Remark    string    `json:"remark"`
	Url       string    `json:"url"`
	UserId    uint      `json:"user_id"`
	Token     string    `json:"token"`
}
