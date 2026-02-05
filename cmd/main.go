package main

import (
	"fmt"
	"jing-sync/boot/app"
	"jing-sync/boot/config"
	"jing-sync/boot/database"
	"jing-sync/boot/i18n"
	"jing-sync/boot/logger"
	"jing-sync/internal/services"
	"jing-sync/internal/services/db_services"

	"github.com/robfig/cron/v3"
	"time"
)

func main() {
	// 初始化配置
	config.InitConfig(config.EnvTypeProd)
	// 初始化日志
	logger.LoggerInit()
	defer logger.GetLogger().Writer().Close()
	i18n.I18nInit()
	// 初始化数据库
	password := database.InitDB()

	// 启动cron调度器
	c := cron.New()

	// 加载并调度所有启用的任务
	scheduleJobs(c)

	c.Start()
	defer c.Stop()

	// web设置
	r := app.WebSet()

	if password != "" {
		logger.GetLogger().Infof("admin password:%s", password)
	}

	port := config.Cfg.Port
	logger.GetLogger().Infof("^_^ Running at http://127.0.0.1:%d/", port)
	r.Run(fmt.Sprintf(":%d", port))
}

// scheduleJobs 加载所有任务并设置cron调度
func scheduleJobs(c *cron.Cron) {
	db := database.GetDB()
	jobService := db_services.NewJobService(db)

	// 获取所有启用的任务
	jobs, _ := jobService.GetPageList(1, 100)
	if jobs == nil {
		return
	}

	for _, job := range jobs.List {
		// 根据调度类型设置cron表达式
		var cronExpr string
		if job.IsCron == 1 {
			// 使用自定义cron表达式 (格式: 分 时 日 月 周)
			cronExpr = fmt.Sprintf("%d %d * * *", job.Interval, job.Interval)
		} else {
			// 使用间隔调度 (Interval单位为分钟)
			cronExpr = fmt.Sprintf("*/%d * * * *", job.Interval)
		}

		// 添加调度任务
		jobId := job.ID
		_, err := c.AddFunc(cronExpr, func() {
			logger.GetLogger().Infof("开始执行同步任务: %s (ID: %d)", job.Remark, jobId)
			startTime := time.Now()

			syncService := services.NewSyncService(db)
			syncLog, err := syncService.ExecuteSync(jobId)

			if err != nil {
				logger.GetLogger().Errorf("同步任务执行失败: %s (ID: %d), 错误: %v", job.Remark, jobId, err)
			} else if syncLog.Status == 1 {
				logger.GetLogger().Infof("同步任务执行成功: %s (ID: %d), 耗时: %v, 新增文件: %d",
					job.Remark, jobId, time.Since(startTime), syncLog.FilesAdded)
			} else {
				logger.GetLogger().Errorf("同步任务部分失败: %s (ID: %d), 错误: %s",
					job.Remark, jobId, syncLog.ErrorMsg)
			}
		})

		if err != nil {
			logger.GetLogger().Errorf("任务调度失败: %s (ID: %d), 表达式: %s, 错误: %v",
				job.Remark, job.ID, cronExpr, err)
		} else {
			logger.GetLogger().Infof("任务已调度: %s (ID: %d), 表达式: %s",
				job.Remark, job.ID, cronExpr)
		}
	}
}
