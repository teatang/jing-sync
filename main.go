package main

import (
	"fmt"
	"jing-sync/boot"
	"jing-sync/config"
	"jing-sync/logger"
)

func main() {
	// 初始化配置
	config.InitConfig()
	// 初始化日志
	logger.LoggerInit()
	defer logger.Log.Writer().Close()

	// 初始化数据库
	password := boot.InitDB()

	// web设置
	r := boot.WebSet()

	if password != "" {
		logger.Log.Info(fmt.Sprintf("admin password:%s", password))
	}

	port := config.Cfg.Port
	logger.Log.Info(fmt.Sprintf("^_^ Running at http://127.0.0.1:%d/", port))
	r.Run(fmt.Sprintf(":%d", port))
}
