package models

import "time"

type Job struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    uint8     `json:"status" gorm:"default:1"`
	Remark    string    `json:"remark"`
	SrcPath   string    `json:"src_path"`
	DstPath   string    `json:"dst_path"`
	EngineId  uint      `json:"engine_id"`
	Speed     uint8     `json:"speed"`
	Method    string    `json:"method"`
	Interval  uint      `json:"interval"`
	IsCron    uint8     `json:"is_cron"`
}
