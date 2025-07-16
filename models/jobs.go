package models

type Job struct {
	BaseModel
	Remark    string    `json:"remark"`
	SrcPath   string    `json:"src_path"`
	DstPath   string    `json:"dst_path"`
	EngineId  uint      `json:"engine_id"`
	Speed     uint8     `json:"speed"`
	Method    string    `json:"method"`
	Interval  uint      `json:"interval"`
	IsCron    uint8     `json:"is_cron"`
}
