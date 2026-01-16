package llm

import (
	"context"
	"fmt"
	"sk-im-bot/internal/config"

	openai "github.com/sashabaranov/go-openai"
)

type LLMClient struct {
	client *openai.Client
	config config.LLMConfig
}

func NewLLMClient(cfg config.LLMConfig) *LLMClient {
	openaiConfig := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		openaiConfig.BaseURL = cfg.BaseURL
	}

	client := openai.NewClientWithConfig(openaiConfig)
	return &LLMClient{
		client: client,
		config: cfg,
	}
}

func (l *LLMClient) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := l.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     l.config.Model,
			Messages:  messages,
			MaxTokens: l.config.MaxTokens,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return resp.Choices[0].Message.Content, nil
}
