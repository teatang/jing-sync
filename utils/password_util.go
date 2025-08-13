package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"os"
	"github.com/golang-jwt/jwt/v5"
)

// JWT声明结构
type Claims struct {
	Username string `json:"username"`
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Password2hash(password string) (string, error) {
	sk, err := GetSecretKey()
	if err != nil {
		return "", err
	}
	s := password + sk
	return SHA256(s), nil
}

// 使用crypto/rand生成安全随机字符串
func SecureRandString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[num.Int64()]
	}
	return string(b), nil
}

func GetSecretKey() (string, error) {
	secretKey, _ := SecureRandString(256)
	return ReadOrSet("data/secret.key", secretKey, false)
}

// ReadOrSet 读取文件内容，不存在则创建
// fileName: 文件路径
// defaultValue: 默认值
// force: 是否强制覆盖
func ReadOrSet(fileName string, defaultValue string, force bool) (string, error) {
	if !force {
		if _, err := os.Stat(fileName); err == nil {
			data, err := os.ReadFile(fileName)
			if err != nil {
				return "", err
			}
			return string(data), nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}

	err := os.WriteFile(fileName, []byte(defaultValue), 0644)
	if err != nil {
		return "", err
	}
	return defaultValue, nil
}

// SHA256加密
func SHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
