package tests

import (
	"fmt"
	"jing-sync/boot"
	"jing-sync/services/db_services"
	"jing-sync/utils"

	"testing"
	"unicode/utf8"
)

// 测试RandString基础功能
func TestRandString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"常规长度10", 10},
		{"常规长度20", 20},
		{"边界长度0", 0},
		{"边界长度1", 1},
		{"大长度测试", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := utils.SecureRandString(tt.size)

			// 验证长度
			if utf8.RuneCountInString(got) != tt.size {
				t.Errorf("生成字符串：%s，期望长度 %d, 实际长度 %d", got, tt.size, len(got))
			}

			// 验证字符范围（仅当长度>0时）
			if tt.size > 0 {
				for _, r := range got {
					if !isValidRune(r) {
						t.Errorf("非法字符 %U", r)
					}
				}
			}
		})
	}
}

// 测试SecureRandString安全随机性
func TestSecureRandString(t *testing.T) {
	const testSize = 10
	str1, err := utils.SecureRandString(testSize)

	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	str2, err := utils.SecureRandString(testSize)
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	t.Log(str1, str2)
	// 验证两次生成结果不同（概率性测试）
	if str1 == str2 {
		t.Errorf("两次生成相同结果，不符合随机性要求")
	}
}

// 辅助函数：验证字符是否在允许范围内
func isValidRune(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9')
}

func TestGetUserByUsernamePassword(t *testing.T) {
	db := boot.GetDB()
    us := db_services.NewUserService(db)

	user, _ := us.GetUserByUsernamePassword("admin", "78rzLF8uXg")

	fmt.Println(user)
}
