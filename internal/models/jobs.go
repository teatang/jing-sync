package models

// 同步模式常量定义
const (
	SyncModeFull    uint8 = 0 // 全同步：让目标目录与源目录完全一致
	SyncModeAddOnly uint8 = 1 // 仅新增：目标目录缺少的文件进行复制，不覆盖不删除
	SyncModeMove    uint8 = 2 // 移动模式：同步后删除源文件
)

// 同步速度/频率定义
const (
	SyncSpeedNormal uint8 = 0 // 正常速度
	SyncSpeedFast   uint8 = 1 // 快速（不强制刷新目录缓存）
	SyncSpeedSlow   uint8 = 2 // 慢速（操作间隔 1 秒）
)

type Job struct {
	BaseModel
	Remark   string `json:"remark" gorm:"column:remark;size:255"`
	SrcPath  string `json:"src_path" gorm:"column:src_path;size:500"`
	DstPath  string `json:"dst_path" gorm:"column:dst_path;size:500"`
	EngineId uint   `json:"engine_id" gorm:"column:engine_id"`
	UserId   uint   `json:"user_id" gorm:"column:user_id"`
	Speed    uint8  `json:"speed" gorm:"column:speed;default:0"`       // 速度策略
	SyncMode uint8  `json:"sync_mode" gorm:"column:sync_mode;default:1"` // 0:全同步 1:仅新增 2:移动
	Method   string `json:"method" gorm:"column:method;size:50"`
	Interval uint   `json:"interval" gorm:"column:interval;default:60"`   // 间隔分钟数
	IsCron   uint8  `json:"is_cron" gorm:"column:is_cron;default:0"`     // 0:间隔 1:Cron表达式
	CronExpr string `json:"cron_expr" gorm:"column:cron_expr;size:50"`  // Cron 表达式
	MinSize  *int64 `json:"min_size" gorm:"column:min_size"`             // 过滤：最小文件字节
	MaxSize  *int64 `json:"max_size" gorm:"column:max_size"`             // 过滤：最大文件字节
	Excludes string `json:"excludes" gorm:"column:excludes;type:text"`   // 换行分隔的排除规则
	Enabled  uint8  `json:"enabled" gorm:"column:enabled;default:1"`     // 0:禁用 1:启用
}
