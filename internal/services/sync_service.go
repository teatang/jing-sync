package services

import (
	"encoding/json"
	"fmt"
	"jing-sync/internal/models"
	"jing-sync/internal/services/db_services"
	"jing-sync/internal/services/filter"
	"jing-sync/internal/utils"
	"time"

	"gorm.io/gorm"
)

type SyncService struct {
	db            *gorm.DB
	jobService    *db_services.JobService
	engineService *db_services.EngineService
}

func NewSyncService(db *gorm.DB) *SyncService {
	return &SyncService{
		db:            db,
		jobService:    db_services.NewJobService(db),
		engineService: db_services.NewEngineService(db),
	}
}

// SyncResult 记录同步详细统计与过程
type SyncResult struct {
	JobId        uint     `json:"job_id"`
	JobRemark    string   `json:"job_remark"`
	SrcPath      string   `json:"src_path"`
	DstPath      string   `json:"dst_path"`
	EngineUrl    string   `json:"engine_url"`
	SyncMode     uint8    `json:"sync_mode"`
	FilesFound   int      `json:"files_found"`
	FilesSkipped int      `json:"files_skipped"`
	FilesAdded   int      `json:"files_added"`
	FilesFailed  int      `json:"files_failed"`
	DirsFound    int      `json:"dirs_found"`
	Details      []string `json:"details"`
}

// ExecuteSync 执行同步任务全流程
func (s *SyncService) ExecuteSync(jobId uint) (*models.SyncLog, error) {
	startTime := time.Now()

	// 1. 获取任务信息
	job, err := s.jobService.GetByID(utils.UintToString(jobId))
	if err != nil {
		return s.createSyncLog(jobId, startTime, nil, "获取作业配置失败: "+err.Error())
	}

	// 2. 检查源路径与目的路径重叠情况（安全保护，防止死循环）
	if filter.VirtualPathsOverlap(job.SrcPath, job.DstPath) {
		errOverlap := fmt.Errorf("source and destination paths overlap: %s and %s", job.SrcPath, job.DstPath)
		return s.createSyncLog(jobId, startTime, nil, errOverlap.Error())
	}

	// 3. 获取存储引擎配置
	engine, err := s.engineService.GetByID(utils.UintToString(job.EngineId))
	if err != nil {
		return s.createSyncLog(jobId, startTime, nil, "获取引擎配置失败: "+err.Error())
	}

	// 4. 执行实际比对与同步传输
	result, err := s.doSync(job, engine)

	// 5. 记录并持久化同步日志
	return s.createSyncLogWithResult(jobId, startTime, result, err)
}

// doSync 执行真正的扫描、过滤、差异比对与数据同步
func (s *SyncService) doSync(job *models.Job, engine *models.Engine) (*SyncResult, error) {
	result := &SyncResult{
		JobId:     job.ID,
		JobRemark: job.Remark,
		SrcPath:   job.SrcPath,
		DstPath:   job.DstPath,
		EngineUrl: engine.Url,
		SyncMode:  job.SyncMode,
		Details:   make([]string, 0),
	}

	client := NewOpenListClient(utils.UintToString(engine.ID), s.db)
	pathFilter := filter.NewPathFilter(job.Excludes)

	// 1. 扫描源端目录文件
	result.Details = append(result.Details, fmt.Sprintf("开始扫描源端目录: %s", job.SrcPath))
	srcFiles, err := client.ListFilesDetailed(job.SrcPath, int(job.Speed))
	if err != nil {
		errMsg := fmt.Sprintf("扫描源目录失败: %v", err)
		result.Details = append(result.Details, errMsg)
		return result, err
	}

	// 2. 如果是“仅新增”模式，先扫描目标目录获取已有文件清单
	dstFileMap := make(map[string]FileInfo)
	if job.SyncMode == models.SyncModeAddOnly || job.SyncMode == models.SyncModeFull {
		dstFiles, err := client.ListFilesDetailed(job.DstPath, int(job.Speed))
		if err == nil {
			for _, df := range dstFiles {
				if !df.IsDir {
					dstFileMap[df.Name] = df
				}
			}
		}
	}

	// 3. 逐个遍历源端文件，执行过滤与差异同步
	for _, sf := range srcFiles {
		if sf.IsDir {
			result.DirsFound++
			continue
		}
		result.FilesFound++

		// 检查路径/名称排除通配符
		if pathFilter.IsExcluded(sf.Name) {
			result.FilesSkipped++
			result.Details = append(result.Details, fmt.Sprintf("跳过文件 (规则排除): %s", sf.Name))
			continue
		}

		// 检查文件大小限制
		if !filter.IsFileSizeAllowed(sf.Size, job.MinSize, job.MaxSize) {
			result.FilesSkipped++
			result.Details = append(result.Details, fmt.Sprintf("跳过文件 (大小超出限制 %s): %s", utils.FormatBytes(sf.Size), sf.Name))
			continue
		}

		// 增量模式：如果目标端已存在相同文件且大小一致，则跳过
		if job.SyncMode == models.SyncModeAddOnly {
			if df, exists := dstFileMap[sf.Name]; exists {
				if df.Size == sf.Size {
					result.FilesSkipped++
					result.Details = append(result.Details, fmt.Sprintf("跳过文件 (目标已存在): %s", sf.Name))
					continue
				}
			}
		}

		// 4. 调用 OpenList API 触发同步或移动
		var opErr error
		var taskId string

		if job.SyncMode == models.SyncModeMove {
			taskId, opErr = client.MoveFile(job.SrcPath, job.DstPath, sf.Name)
		} else {
			taskId, opErr = client.CopyFile(job.SrcPath, job.DstPath, sf.Name)
		}

		if opErr != nil {
			result.FilesFailed++
			result.Details = append(result.Details, fmt.Sprintf("同步失败: %s, 错误: %v", sf.Name, opErr))
		} else {
			result.FilesAdded++
			result.Details = append(result.Details, fmt.Sprintf("同步成功: %s, 任务ID: %s", sf.Name, taskId))
		}
	}

	result.Details = append(result.Details, fmt.Sprintf("同步完成，发现文件: %d, 新增: %d, 跳过: %d, 失败: %d",
		result.FilesFound, result.FilesAdded, result.FilesSkipped, result.FilesFailed))

	return result, nil
}

// createSyncLog 持久化同步日志
func (s *SyncService) createSyncLog(jobId uint, startTime time.Time, result *SyncResult, errMsg string) (*models.SyncLog, error) {
	status := models.SyncStatusSuccess
	errorMsg := ""
	endTime := time.Now()

	if errMsg != "" {
		status = models.SyncStatusFailed
		errorMsg = errMsg
	}

	var details string
	if result != nil {
		detailsBytes, _ := json.Marshal(result)
		details = string(detailsBytes)
	}

	filesAdded := 0
	if result != nil {
		filesAdded = result.FilesAdded
	}

	syncLog := &models.SyncLog{
		JobId:      jobId,
		Status:     status,
		ErrorMsg:   errorMsg,
		FilesAdded: filesAdded,
		StartTime:  startTime,
		EndTime:    endTime,
		Details:    details,
	}

	// 保存到数据库
	err := s.db.Create(syncLog).Error
	if errMsg != "" {
		return syncLog, fmt.Errorf("%s", errMsg)
	}
	return syncLog, err
}

func (s *SyncService) createSyncLogWithResult(jobId uint, startTime time.Time, result *SyncResult, err error) (*models.SyncLog, error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return s.createSyncLog(jobId, startTime, result, errMsg)
}

// GetSyncLogsByJobId 分页获取任务的历史同步日志
func (s *SyncService) GetSyncLogsByJobId(jobId uint, page, size int) (*utils.PageList[models.SyncLog], error) {
	var logs []models.SyncLog
	offset := (page - 1) * size

	query := s.db.Where("job_id = ?", jobId).Order("id DESC")
	err := query.Offset(offset).Limit(size).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	var count int64
	s.db.Model(&models.SyncLog{}).Where("job_id = ?", jobId).Count(&count)

	return &utils.PageList[models.SyncLog]{
		List: logs,
		Pagination: utils.PageInfo{
			Page:  page,
			Size:  size,
			Total: count,
		},
	}, nil
}

// GetSyncLogById 获取单条日志详细记录
func (s *SyncService) GetSyncLogById(id uint) (*models.SyncLog, error) {
	var log models.SyncLog
	err := s.db.First(&log, id).Error
	return &log, err
}

// ListAllSyncLogs 获取全局最新同步日志
func (s *SyncService) ListAllSyncLogs(page, size int) (*utils.PageList[models.SyncLog], error) {
	var logs []models.SyncLog
	offset := (page - 1) * size

	err := s.db.Order("id DESC").Offset(offset).Limit(size).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	var count int64
	s.db.Model(&models.SyncLog{}).Count(&count)

	return &utils.PageList[models.SyncLog]{
		List: logs,
		Pagination: utils.PageInfo{
			Page:  page,
			Size:  size,
			Total: count,
		},
	}, nil
}
