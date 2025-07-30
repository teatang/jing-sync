package tests

import (
	"fmt"
	"jing-sync/boot"
	"jing-sync/config"
	"jing-sync/services"
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
		fmt.Printf("Failed to get current working directory: %s", err)
	}

	// 这是单元测试
	if strings.HasSuffix(rootDir, "tests") {
		rootDir = filepath.Dir(rootDir)
	}

	// 设置工作目录
	if err := os.Chdir(rootDir); err != nil {
		fmt.Printf("Failed to change working directory: %s", err)
	}
	fmt.Println("工作目录：", rootDir)
}

func TestOpenListClientPost(t *testing.T) {
	ol := services.NewOpenListClient("1", boot.GetDB())
	data := map[string]interface{}{
		"path":     "/",
		"refresh":  true,
		"page":     1,
		"per_page": 100,
	}

	res, _ := ol.Post("/api/fs/list", data)

	fmt.Println(string(res))
}
