package main

import (
	"jing-sync/boot"
	"jing-sync/logger"
)

func main() {
	// 初始化日志
	logger.LoggerInit()
	defer logger.Log.Writer().Close()

	// 初始化数据库
	boot.InitDB()

	// 初始化路由
	r := boot.SetupRouter()

	logger.Log.Info("服务启动成功 Port:8888")
	r.Run(":8888")
}
