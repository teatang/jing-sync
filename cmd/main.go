package main

import (
	"fmt"
	"jing-sync/boot/app"
	"jing-sync/boot/config"
	"jing-sync/boot/database"
	"jing-sync/boot/i18n"
	"jing-sync/boot/logger"
	"jing-sync/internal/services"
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
	db := database.GetDB()

	// 启动全局 Cron 调度管理器
	services.GetCronManager(db)

	// web设置
	r := app.WebSet()

	if password != "" {
		logger.GetLogger().Infof("admin password:%s", password)
	}

	port := config.Cfg.Port
	logger.GetLogger().Infof("^_^ Running at http://127.0.0.1:%d/", port)
	r.Run(fmt.Sprintf(":%d", port))
}
