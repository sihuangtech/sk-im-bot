package llm

import (
	"context"
	"fmt"
	"sk-im-bot/internal/config"

	openai "github.com/sashabaranov/go-openai"
)

// LLMClient 封装了大模型对话接口，底层由 OpenAI SDK 驱动
type LLMClient struct {
	client *openai.Client   // 通用的 API 请求实例
	config config.LLMConfig // 合并加载的动态配置项
}

// NewLLMClient 注入配置并初始化相应的网络请求凭据及 BaseURL 地址
func NewLLMClient(cfg config.LLMConfig) *LLMClient {
	// 使用 SDK 默认配置，注入用户提供的 API Key
	openaiConfig := openai.DefaultConfig(cfg.APIKey)

	// 如果配置了自定义网关（代理地址），则修改默认的 API 端点
	if cfg.BaseURL != "" {
		openaiConfig.BaseURL = cfg.BaseURL
	}

	// 实例化 SDK 客户端句柄
	client := openai.NewClientWithConfig(openaiConfig)
	return &LLMClient{
		client: client,
		config: cfg,
	}
}

// Chat 发起一次聊天补全请求。messages 参数支持历史会话传入，从而实现多轮对话
func (l *LLMClient) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	// 构造 OpenAI 规范格式的对话生成请求
	resp, err := l.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     l.config.Model,     // 使用配置中定义的模型版本
			Messages:  messages,           // 对话上下文堆栈
			MaxTokens: l.config.MaxTokens, // 生成回复的最大词数限制
		},
	)

	// 网络请求或 API 层的错误处理
	if err != nil {
		return "", fmt.Errorf("LLM API 服务器异常: %w", err)
	}

	// 边界检查：如果 API 返回结果集为空，则抛出逻辑错误
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("LLM 提供商未返回有效的对话响应内容")
	}

	// 提取生成候选集中的第一条内容作为回复
	return resp.Choices[0].Message.Content, nil
}
