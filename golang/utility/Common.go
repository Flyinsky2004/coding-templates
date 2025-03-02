package utility

import (
	"math/rand"
	"strings"
	"time"
)

// GenerateCode 生成随机验证码
func GenerateCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}
func GenerateRandomString(n int) string {
	// 定义随机字符集合
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 使用strings.Builder优化字符串拼接
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}
	return sb.String()
}
