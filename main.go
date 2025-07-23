package main

import (
	"fmt"
	"jing-sync/boot"
	"jing-sync/logger"
)

func main() {
	// 初始化日志
	logger.LoggerInit()
	defer logger.Log.Writer().Close()

	// 初始化数据库
	password := boot.InitDB()

	// 初始化路由
	r := boot.SetupRouter()

	if password != "" {
		logger.Log.Info(fmt.Sprintf("admin password:%s", password))
	}
	logger.Log.Info("服务启动成功 Port:8888")
	r.Run(":8888")
}
