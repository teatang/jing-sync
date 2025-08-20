package main

import (
	"fmt"
	"jing-sync/boot/app"
	"jing-sync/boot/config"
	"jing-sync/boot/database"
	"jing-sync/boot/logger"

	"github.com/robfig/cron/v3"
	"time"
)

func main() {
	// 初始化配置
	config.InitConfig(config.EnvTypeProd)
	// 初始化日志
	logger.LoggerInit()
	defer logger.GetLogger().Writer().Close()

	// 初始化数据库
	password := database.InitDB()

	c := cron.New()
	i := 1
	EntryID, err := c.AddFunc("*/1 * * * *", func() {
		fmt.Println(time.Now(), "每分钟执行一次", i)
		i++
	})
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
	defer c.Stop()

	// web设置
	r := app.WebSet()

	if password != "" {
		logger.GetLogger().Info(fmt.Sprintf("admin password:%s", password))
	}

	port := config.Cfg.Port
	logger.GetLogger().Info(fmt.Sprintf("^_^ Running at http://127.0.0.1:%d/", port))
	r.Run(fmt.Sprintf(":%d", port))
}
