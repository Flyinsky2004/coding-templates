package utility

import (
	"bytes"
	"example/config"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	baseURL      = "https://%s.tts.speech.microsoft.com/cognitiveservices/v1"
	ssmlTemplate = `
<speak version='1.0' xml:lang='%s'>
    <voice xml:lang='%s' xml:gender='%s' name='%s'>
        %s
    </voice>
</speak>`
	basePath = "/Users/wangjiying/Documents/recording/legacy/audio"
)

type AzureTTS struct {
	SubscriptionKey string
	Region          string
	client          *http.Client
}

// NewAzureTTS 创建新的Azure TTS客户端
func NewAzureTTS(subscriptionKey, region string) *AzureTTS {
	// 如果未提供参数，则使用配置中的默认值
	if subscriptionKey == "" {
		subscriptionKey = config.Config.AzureTTS.SubscriptionKey
	}
	if region == "" {
		region = config.Config.AzureTTS.Region
	}

	return &AzureTTS{
		SubscriptionKey: subscriptionKey,
		Region:          region,
		client:          &http.Client{Timeout: 180 * time.Second},
	}
}

// TextToSpeech 将文本转换为语音并保存到文件
func (tts *AzureTTS) TextToSpeech(text, lang, voice, gender string) (string, error) {
	// 生成文件名
	randomString := func() string {
		const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
		b := make([]byte, 6)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}
	timestamp := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s.mp3", timestamp, randomString())
	outputPath := fmt.Sprintf(basePath+"/%s", filename)

	// 确保目录存在
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 检查文本长度并分段
	const maxLength = 2000
	textRunes := []rune(text)
	if len(textRunes) > maxLength {
		// 创建临时目录
		tempDir := basePath + "/temp"
		if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("创建临时目录失败: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// 分段处理
		var tempFiles []string
		for i := 0; i < len(textRunes); i += maxLength {
			end := i + maxLength
			if end > len(textRunes) {
				end = len(textRunes)
			}
			chunk := string(textRunes[i:end])

			// 为每个分段生成临时文件名
			tempFile := fmt.Sprintf("%s/chunk_%d.mp3", tempDir, i/maxLength)
			tempFiles = append(tempFiles, tempFile)

			// 处理当前分段
			if err := tts.processChunk(chunk, lang, voice, gender, tempFile); err != nil {
				return "", fmt.Errorf("处理文本分段失败: %v", err)
			}
		}

		// 合并所有临时文件
		if err := tts.mergeAudioFiles(tempFiles, outputPath); err != nil {
			return "", fmt.Errorf("合并音频文件失败: %v", err)
		}

		return filename, nil
	}

	// 文本长度在限制内，直接处理
	if err := tts.processChunk(text, lang, voice, gender, outputPath); err != nil {
		return "", err
	}

	return filename, nil
}

// processChunk 处理单个文本块
func (tts *AzureTTS) processChunk(text, lang, voice, gender, outputPath string) error {
	// 构建SSML
	ssml := fmt.Sprintf(config.Config.AzureTTS.SSMLTemplate, lang, lang, gender, voice, text)

	// 创建请求
	url := fmt.Sprintf(config.Config.AzureTTS.BaseURL, tts.Region)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(ssml))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", tts.SubscriptionKey)
	req.Header.Set("X-Microsoft-OutputFormat", "audio-16khz-128kbitrate-mono-mp3")

	// 发送请求
	resp, err := tts.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API调用失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 写入文件
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// mergeAudioFiles 合并多个音频文件
func (tts *AzureTTS) mergeAudioFiles(tempFiles []string, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建合并文件失败: %v", err)
	}
	defer outFile.Close()

	for _, tempFile := range tempFiles {
		data, err := os.ReadFile(tempFile)
		if err != nil {
			return fmt.Errorf("读取临时文件失败: %v", err)
		}
		if _, err := outFile.Write(data); err != nil {
			return fmt.Errorf("写入合并文件失败: %v", err)
		}
	}

	return nil
}
