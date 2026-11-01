package services

import (
	"fmt"
	"jing-sync/boot/logger"
	"jing-sync/internal/models"
	"jing-sync/internal/services/db_services"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronManager struct {
	mu      sync.Mutex
	cron    *cron.Cron
	db      *gorm.DB
	entries map[uint]cron.EntryID
}

var globalCronManager *CronManager
var once sync.Once

func GetCronManager(db *gorm.DB) *CronManager {
	once.Do(func() {
		globalCronManager = &CronManager{
			cron:    cron.New(),
			db:      db,
			entries: make(map[uint]cron.EntryID),
		}
		globalCronManager.Start()
	})
	return globalCronManager
}

func (cm *CronManager) Start() {
	cm.cron.Start()
	cm.ReloadAll()
}

func (cm *CronManager) Stop() {
	cm.cron.Stop()
}

func (cm *CronManager) ReloadAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 清除现有的所有调度
	for jobId, entryId := range cm.entries {
		cm.cron.Remove(entryId)
		delete(cm.entries, jobId)
	}

	jobService := db_services.NewJobService(cm.db)
	jobs, err := jobService.GetPageList(1, 1000)
	if err != nil || jobs == nil {
		return
	}

	for _, job := range jobs.List {
		cm.scheduleJobLocked(&job)
	}
}

func (cm *CronManager) AddOrUpdateJob(job *models.Job) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 先移除旧的调度（如果存在），防止重复添加
	if entryId, exists := cm.entries[job.ID]; exists {
		cm.cron.Remove(entryId)
		delete(cm.entries, job.ID)
		logger.GetLogger().Infof("已移除旧的任务调度 (更新): Job ID %d", job.ID)
	}

	// 如果任务已启用且不是仅手动触发 (is_cron != 2)
	if job.Enabled == 1 && job.IsCron != 2 {
		cm.scheduleJobLocked(job)
	}
}

func (cm *CronManager) RemoveJob(jobId uint) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if entryId, exists := cm.entries[jobId]; exists {
		cm.cron.Remove(entryId)
		delete(cm.entries, jobId)
		logger.GetLogger().Infof("已移除任务调度 (删除): Job ID %d", jobId)
	}
}

func (cm *CronManager) scheduleJobLocked(job *models.Job) {
	if job.IsCron == 2 || job.Enabled == 0 {
		return
	}

	var cronExpr string
	if job.IsCron == 1 {
		cronExpr = job.CronExpr
		if cronExpr == "" {
			cronExpr = "0 2 * * *"
		}
	} else {
		interval := job.Interval
		if interval == 0 {
			interval = 60
		}
		cronExpr = fmt.Sprintf("*/%d * * * *", interval)
	}

	jobId := job.ID
	jobRemark := job.Remark
	db := cm.db

	entryId, err := cm.cron.AddFunc(cronExpr, func() {
		logger.GetLogger().Infof("开始执行同步任务: %s (ID: %d)", jobRemark, jobId)
		startTime := time.Now()

		syncService := NewSyncService(db)
		syncLog, err := syncService.ExecuteSync(jobId)

		if err != nil {
			logger.GetLogger().Errorf("同步任务执行失败: %s (ID: %d), 错误: %v", jobRemark, jobId, err)
		} else if syncLog.Status == 1 {
			logger.GetLogger().Infof("同步任务执行成功: %s (ID: %d), 耗时: %v, 新增文件: %d",
				jobRemark, jobId, time.Since(startTime), syncLog.FilesAdded)
		} else {
			logger.GetLogger().Errorf("同步任务部分失败: %s (ID: %d), 错误: %s",
				jobRemark, jobId, syncLog.ErrorMsg)
		}
	})

	if err != nil {
		logger.GetLogger().Errorf("任务调度失败: %s (ID: %d), 表达式: %s, 错误: %v",
			jobRemark, jobId, cronExpr, err)
	} else {
		cm.entries[jobId] = entryId
		logger.GetLogger().Infof("任务已调度: %s (ID: %d), 表达式: %s",
			jobRemark, jobId, cronExpr)
	}
}
