package utility

import (
	"bytes"
	"context"
	"encoding/json"
	"example/config"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sashabaranov/go-openai"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/29 17:25
 */

// ==================== 消息结构定义 ====================

// ChatRequest 定义发送给 ChatGPT API 的请求结构
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Prompt      string    `json:"prompt,omitempty"`
	Question    string    `json:"question,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

// Message 定义消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 定义 ChatGPT API 返回的响应结构
type ChatMetaResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// StreamResponse 定义流式响应结构
type StreamResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// ==================== 图像生成相关 ====================

// Dalle3Request 表示 DALL-E 3 API 的请求体
type Dalle3Request struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
	Model  string `json:"model"`
}

// Dalle3Response 表示 DALL-E 3 API 的响应
type Dalle3MetaResponse struct {
	Data []struct {
		URL            string `json:"url"`
		Revised_Prompt string `json:"revised_prompt"`
	} `json:"data"`
}

// GenerateImage 使用 DALL-E 3 API 生成图像
func GenerateImage(prompt string) (string, error) {
	// 使用配置中的基础URL和API密钥
	baseURL := config.Config.OpenAI.BaseURL
	apiKey := config.Config.OpenAI.Key

	// 构造请求体
	requestBody := Dalle3Request{
		Prompt: prompt,
		N:      1,           // 生成图像的数量
		Size:   "1024x1024", // 图像大小
		Model:  "dall-e-3",
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/v1/images/generations", baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// 读取并解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var dalleResponse Dalle3MetaResponse
	err = json.Unmarshal(body, &dalleResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// 检查响应中是否至少有一个URL
	if len(dalleResponse.Data) == 0 {
		return "", fmt.Errorf("no image URL returned in response")
	}

	// 返回第一个图像URL
	return dalleResponse.Data[0].URL, nil
}

// DownloadImage 下载图像并保存到本地
func DownloadImage(imageURL string) (string, error) {
	// 创建 /uploads 目录（如果不存在）
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// 生成唯一文件名
	rand.Seed(time.Now().UnixNano())
	randomSuffix := make([]byte, 3) // 3 bytes = 6 hex characters
	if _, err := rand.Read(randomSuffix); err != nil {
		return "", fmt.Errorf("failed to generate random suffix: %v", err)
	}
	fileName := fmt.Sprintf("%s_%06x.webp", time.Now().Format("20060102"), randomSuffix)
	filePath := filepath.Join(uploadDir, fileName)

	// 下载图像
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	// 创建输出文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// 将图像数据写入文件
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	return fileName, nil
}

// ==================== OpenAI 客户端 ====================

// Client 使用 OpenAI SDK 客户端
type Client struct {
	client  *openai.Client
	baseURL string
}

// NewClient 创建一个新的客户端
func NewClient(apiKey string) *Client {
	openaiConfig := openai.DefaultConfig(apiKey)
	openaiConfig.BaseURL = config.Config.OpenAI.BaseURL

	return &Client{
		client:  openai.NewClientWithConfig(openaiConfig),
		baseURL: config.Config.OpenAI.BaseURL,
	}
}

// SendMessage 使用 OpenAI SDK 发送消息
func (c *Client) SendMessage(messages []Message, model string, maxToken int, temperature float32) (ChatMetaResponse, error) {
	// 转换消息格式
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 创建请求
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       model,
			Messages:    openaiMessages,
			MaxTokens:   maxToken,
			Temperature: float32(temperature),
		},
	)
	if err != nil {
		return ChatMetaResponse{}, fmt.Errorf("failed to create chat completion: %v", err)
	}

	// 转换响应格式为原有的 ChatResponse 结构
	chatResp := ChatMetaResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{
				Message: struct {
					Content string `json:"content"`
				}{
					Content: resp.Choices[0].Message.Content,
				},
			},
		},
		Usage: struct {
			TotalTokens int `json:"total_tokens"`
		}{
			TotalTokens: resp.Usage.TotalTokens,
		},
	}

	return chatResp, nil
}

// ChatHandler 保持原有函数签名和行为不变
func ChatHandler(request ChatRequest) (ChatMetaResponse, error) {
	client := NewClient(config.Config.OpenAI.Key)

	systemMessage := []Message{
		{
			Role:    "system",
			Content: request.Prompt,
		},
	}
	askMessage := []Message{
		{
			Role:    "user",
			Content: request.Question,
		},
	}
	messages := append(systemMessage, request.Messages...)
	messages = append(messages, askMessage...)

	response, err := client.SendMessage(messages, request.Model, request.MaxTokens, request.Temperature)
	if err != nil {
		fmt.Println("Error:", err)
		return ChatMetaResponse{}, err
	}

	return response, nil
}

// ==================== 流式聊天 ====================

// StreamChatCompletion 创建流式聊天完成
func StreamChatCompletion(ctx context.Context, request ChatRequest) (<-chan StreamResponse, error) {
	// 创建配置
	openaiConfig := openai.DefaultConfig(config.Config.OpenAI.Key)
	openaiConfig.BaseURL = config.Config.OpenAI.BaseURL

	// 使用配置创建客户端
	client := openai.NewClientWithConfig(openaiConfig)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: request.Prompt,
		},
		{
			Role:    "user",
			Content: request.Question,
		},
	}

	stream, err := client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:       request.Model,
			Messages:    messages,
			Temperature: request.Temperature,
			Stream:      true,
			MaxTokens:   request.MaxTokens,
		},
	)
	if err != nil {
		return nil, err
	}

	responseChan := make(chan StreamResponse)

	go func() {
		defer stream.Close()
		defer close(responseChan)

		for {
			response, err := stream.Recv()
			if err != nil {
				// 流结束
				responseChan <- StreamResponse{
					Done: true,
				}
				return
			}

			if len(response.Choices) > 0 {
				responseChan <- StreamResponse{
					Content: response.Choices[0].Delta.Content,
					Done:    false,
				}
			}
		}
	}()

	return responseChan, nil
}
