package bot

import (
	"sk-im-bot/internal/config"
	"sk-im-bot/pkg/utils"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// DiscordBot 实现针对 Discord 实时的机器人适配器
type DiscordBot struct {
	cfg     config.DiscordConfig // Discord 专用配置 (Token, 频道限制等)
	session *discordgo.Session   // discordgo 的长连接 Session 句柄
	handler func(MessageEvent)   // 收到消息后的分发逻辑处理函数回调
}

// NewDiscordBot 创建一个新的 Discord 机器人适配器实例
func NewDiscordBot(cfg config.DiscordConfig, handler func(MessageEvent)) *DiscordBot {
	return &DiscordBot{
		cfg:     cfg,
		handler: handler,
	}
}

// Start 建立 WebSocket 连接并开始监听 Discord 的 Gateway 事件
func (d *DiscordBot) Start() error {
	var err error
	// 使用 Bot Token 初始化 Session，注意 Token 必须带有 "Bot " 前缀
	d.session, err = discordgo.New("Bot " + d.cfg.Token)
	if err != nil {
		return err
	}

	// 注册消息创建事件处理钩子
	d.session.AddHandler(d.onMessage)

	// 设置 Intent (意图)，确保机器人有权限读取公屏消息和私聊消息
	d.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// 打开长连接
	err = d.session.Open()
	if err != nil {
		utils.Logger.Error("无法开启 Discord 连接", zap.Error(err))
		return err
	}

	utils.Logger.Info("Discord 机器人已上线并处于监听状态")
	return nil
}

// Stop 关闭与 Discord 的连接并释放资源
func (d *DiscordBot) Stop() {
	if d.session != nil {
		d.session.Close()
	}
}

// onMessage 处理来自 Discord 的原始消息推送事件
func (d *DiscordBot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 防死循环：忽略机器人自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 鉴权限制：如果配置了特定服务器 (GuildID)，则忽略不相关的消息
	if d.cfg.GuildID != "" && m.GuildID != d.cfg.GuildID {
		return
	}

	// 判断消息来源是否为群组
	isGroup := m.GuildID != ""

	// 将原始事件包装为统一的内部 MessageEvent 结构并投递给回调
	d.handler(MessageEvent{
		Platform:   "discord",
		PlatformID: m.ChannelID,       // Discord 服务中以频道作为目标
		UserID:     m.Author.ID,       // 发言者的唯一 ID
		Username:   m.Author.Username, // 发言者的昵称
		Content:    m.Content,         // 文本内容
		MsgType:    MsgTypeText,
		IsGroup:    isGroup,
	})
}

// SendMessage 发送一段文本消息到指定的 Discord 频道
func (d *DiscordBot) SendMessage(targetID string, content string, isGroup bool) error {
	_, err := d.session.ChannelMessageSend(targetID, content)
	return err
}

// SendImage 发送一张图片到指定的 Discord 频道 (通过发送链接触发实时预览)
func (d *DiscordBot) SendImage(targetID string, imageUrl string, isGroup bool) error {
	_, err := d.session.ChannelMessageSend(targetID, imageUrl)
	return err
}
