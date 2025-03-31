package utility

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

// SaveBase64ToFile 将base64编码的数据保存为WAV格式文件
func SaveBase64ToFile(base64Data, outputDir, filename string) (string, error) {
	// 确保目录存在
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 移除可能存在的base64头部
	base64Data = strings.TrimPrefix(base64Data, "data:audio/wav;base64,")

	// 解码base64数据
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("解码base64数据失败: %v", err)
	}

	// 确保文件名以.wav结尾
	if !strings.HasSuffix(filename, ".wav") {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".wav"
	}
	tempPath := filepath.Join(outputDir, "temp_"+filename)
	finalPath := filepath.Join(outputDir, filename)

	// 先将base64解码的数据写入临时文件
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return "", fmt.Errorf("写入临时文件失败: %v", err)
	}

	// 使用ffmpeg确保WAV文件为PCM 16-bit、16000Hz、单声道格式
	cmd := exec.Command("ffmpeg", "-i", tempPath,
		"-f", "wav",
		"-acodec", "pcm_s16le", // PCM 16-bit
		"-ar", "16000", // 16000Hz 采样率
		"-ac", "1", // 单声道
		"-y", // 覆盖已存在的文件
		finalPath)

	if err := cmd.Run(); err != nil {
		os.Remove(tempPath) // 清理临时文件
		return "", fmt.Errorf("音频转换失败: %v", err)
	}

	// 删除临时文件
	os.Remove(tempPath)

	return finalPath, nil
}
