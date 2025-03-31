package service

import (
	"example/utility"
	"fmt"
	"math/rand"
	"time"
)

var tts = utility.NewAzureTTS(
	"key content",
	"eastasia", // 你的Azure区域
) // 调用TTS服务

var stt = utility.NewAzureSTT(
	"key content",
	"eastasia", // 你的Azure区域
) // 调用STT服务

func ConvertAudioToText(base64Data string) (string, error) {
	randomString := func() string {
		const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
		b := make([]byte, 6)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}
	timestamp := time.Now().Format("20060102")
	audioFileName := fmt.Sprintf("%s_%s.wav", timestamp, randomString())
	audioPath := "./audio"

	filePath, err := utility.SaveBase64ToFile(base64Data, audioPath, audioFileName)
	answerText, err := stt.SpeechToText(filePath, "zh-CN")
	if err != nil {
		return "nil", err
	}

	return answerText, nil
}

func TextToAudio(content string) (string, error) {
	audioFileName, err := tts.TextToSpeech(content, "zh-CN", "zh-CN-XiaoYanNeural", "Female")
	if err != nil {
		return "", err
	}
	return audioFileName, nil
}
