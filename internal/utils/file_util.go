package utils

import (
	"os"
	"strconv"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// EnsureDir 确保指定的目录存在，如果不存在则创建它。
func EnsureDir(dir string) error {
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(dir, 0755) // MkdirAll会创建所有不存在的父目录
		if err != nil {
			return err
		}
	} else if err != nil {
		// 发生了其他错误（比如权限问题）
		return err
	}
	return nil
}

func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// FormatBytes 将字节大小转换为人类可读的字符串，如 1.00 KB, 2.50 MB
func FormatBytes(bytes int64) string {
	if bytes < 1024 {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	val := float64(bytes)
	for val >= 1024 && i < len(units)-1 {
		val /= 1024
		i++
	}
	return strconv.FormatFloat(val, 'f', 2, 64) + " " + units[i]
}
