package bot

import (
	"context"
	"time"

	"sk-im-bot/internal/config"
	"sk-im-bot/internal/llm"
	"sk-im-bot/internal/model"
	"sk-im-bot/pkg/utils"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// BotManager 核心管理器，协调多个平台的机器人适配器与 LLM 逻辑
type BotManager struct {
	qqAdapter      BotAdapter        // QQ 平台适配器
	discordAdapter BotAdapter        // Discord 平台适配器
	llmClient      *llm.LLMClient    // LLM API 客户端
	msgChan        chan MessageEvent // 全局异步消息处理通道
	// broadcastFunc 用于将收到的实时消息推送到前端 WebSocket
	broadcastFunc func(msg interface{})
}

// Manager 全局机器人管理器单例
var Manager *BotManager

// InitManager 初始化管理器实例及启用的各平台适配器
func InitManager(cfg *config.Config, llmClient *llm.LLMClient, broadcast func(interface{})) {
	Manager = &BotManager{
		llmClient:     llmClient,
		msgChan:       make(chan MessageEvent, 100), // 消息缓冲区容量为 100
		broadcastFunc: broadcast,
	}

	// 根据配置决定是否初始化各平台适配器
	if cfg.QQ.Enabled {
		Manager.qqAdapter = NewQQBot(cfg.QQ, Manager.HandleEvent)
	}
	if cfg.Discord.Enabled {
		Manager.discordAdapter = NewDiscordBot(cfg.Discord, Manager.HandleEvent)
	}
}

// Start 启动所有已激活的适配器并进入主处理循环
func (m *BotManager) Start() {
	if m.qqAdapter != nil {
		go m.qqAdapter.Start()
	}
	if m.discordAdapter != nil {
		go m.discordAdapter.Start()
	}

	// 在独立协程中运行消息分发循环
	go m.processLoop()
}

// HandleEvent 接收来自平台协议层（OneBot/Discordgo）的消息并投递到内部队列
func (m *BotManager) HandleEvent(event MessageEvent) {
	m.msgChan <- event
}

// processLoop 循环处理消息队列中的每一条事件
func (m *BotManager) processLoop() {
	for event := range m.msgChan {
		utils.Logger.Info("处理新消息事件", zap.String("平台", event.Platform), zap.String("内容", event.Content))

		// 1. 异步持久化到数据库
		go m.saveMessage(event)

		// 2. 实时推送到管理控制台
		if m.broadcastFunc != nil {
			m.broadcastFunc(event)
		}

		// 3. 处理机器人自动回复逻辑
		// 备注：此处可扩展判断是否被 @、关键词匹配等
		go m.handleLLMReply(event)
	}
}

// saveMessage 将接收到的消息记录保存到 model 层
func (m *BotManager) saveMessage(event MessageEvent) {
	msg := model.Message{
		Sender:    event.Username,
		Content:   event.Content,
		MsgType:   string(event.MsgType),
		CreatedAt: time.Now(),
	}
	if err := model.DB.Create(&msg).Error; err != nil {
		utils.Logger.Error("消息保存失败", zap.Error(err))
	}
}

// handleLLMReply 调用 LLM 进行对话生成的逻辑入口
func (m *BotManager) handleLLMReply(event MessageEvent) {
	// 【防封号】注入随机等待时间，增加行为仿真度
	RandomDelay(500, 3000)

	// 构造 LLM 对话上下文（此处可根据 Session 获取历史记录进行多轮对话）
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: event.Content},
	}

	response, err := m.llmClient.Chat(context.Background(), messages)
	if err != nil {
		utils.Logger.Error("LLM 接口调用异常", zap.Error(err))
		return
	}

	// 将生成的回答发送回原始平台
	m.SendReply(event.Platform, event.PlatformID, response, event.IsGroup)
}

// SendReply 根据指定的平台分发发送任务
func (m *BotManager) SendReply(platform, targetID, content string, isGroup bool) {
	var err error
	switch platform {
	case "qq":
		if m.qqAdapter != nil {
			err = m.qqAdapter.SendMessage(targetID, content, isGroup)
		}
	case "discord":
		if m.discordAdapter != nil {
			err = m.discordAdapter.SendMessage(targetID, content, isGroup)
		}
	}

	if err != nil {
		utils.Logger.Error("回复发送失败", zap.String("平台", platform), zap.Error(err))
	} else {
		// 成功发送后将机器人自身的回复也存入数据库
		model.DB.Create(&model.Message{
			Sender:    "bot",
			Content:   content,
			MsgType:   "text",
			CreatedAt: time.Now(),
		})
	}
}
