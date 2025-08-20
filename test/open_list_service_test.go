package test

import (
	"jing-sync/boot/config"
	"jing-sync/boot/database"
	"jing-sync/internal/services"

	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	// 初始化配置
	config.InitConfig(config.EnvTypeUnitTest)

	// 获取当前工作目录
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 这是单元测试
	if strings.HasSuffix(rootDir, "tests") {
		rootDir = filepath.Dir(rootDir)
	}

	// 设置工作目录
	if err := os.Chdir(rootDir); err != nil {
		panic(err)
	}
	fmt.Println("工作目录：", rootDir)
}

func TestGetChildPathRaw(t *testing.T) {
	c := services.NewOpenListClient("1", database.GetDB())
	res, _ := c.GetChildPathRaw("/", 0)
	fmt.Println(string(res))
}

func TestGetChildPath(t *testing.T) {
	c := services.NewOpenListClient("1", database.GetDB())
	res, _ := c.GetChildPath("/", 0)
	fmt.Println(res)
}
