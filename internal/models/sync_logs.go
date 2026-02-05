package models

import "time"

type SyncLog struct {
	BaseModel
	JobId      uint   `json:"job_id"`
	Status     uint8  `json:"status"` // 0: 失败, 1: 成功
	ErrorMsg   string `json:"error_msg"`
	FilesAdded int    `json:"files_added"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Details    string `json:"details"` // JSON格式的同步详情
}

const (
	SyncStatusFailed  uint8 = 0
	SyncStatusSuccess uint8 = 1
)
