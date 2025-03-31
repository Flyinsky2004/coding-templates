/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */

package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	sttBaseURL = "https://%s.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1"
)

type AzureSTT struct {
	SubscriptionKey string
	Region          string
	client          *http.Client
}

// NewAzureSTT 创建新的Azure STT客户端
func NewAzureSTT(subscriptionKey, region string) *AzureSTT {
	return &AzureSTT{
		SubscriptionKey: subscriptionKey,
		Region:          region,
		client:          &http.Client{Timeout: 180 * time.Second},
	}
}

type SpeechRecognitionResponse struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	DisplayText       string `json:"DisplayText"`
	Offset            int64  `json:"Offset"`
	Duration          int64  `json:"Duration"`
}

// SpeechToText 将音频文件转换为文本
func (stt *AzureSTT) SpeechToText(audioFilePath string, language string) (string, error) {
	// 读取音频文件
	audioData, err := os.ReadFile(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("读取音频文件失败: %v", err)
	}

	// 构建请求URL
	url := fmt.Sprintf(sttBaseURL, stt.Region)
	if language != "" {
		url += fmt.Sprintf("?language=%s", language)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(audioData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "audio/wav")
	req.Header.Set("Ocp-Apim-Subscription-Key", stt.SubscriptionKey)
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := stt.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API调用失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result SpeechRecognitionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.RecognitionStatus != "Success" {
		return "", fmt.Errorf("语音识别失败: %s", result.RecognitionStatus)
	}

	return result.DisplayText, nil
}
