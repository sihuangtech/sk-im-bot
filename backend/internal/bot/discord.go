package bot

import (
	"sk-im-bot/internal/config"
	"sk-im-bot/pkg/utils"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type DiscordBot struct {
	cfg     config.DiscordConfig
	session *discordgo.Session
	handler func(MessageEvent)
}

func NewDiscordBot(cfg config.DiscordConfig, handler func(MessageEvent)) *DiscordBot {
	return &DiscordBot{
		cfg:     cfg,
		handler: handler,
	}
}

func (d *DiscordBot) Start() error {
	var err error
	d.session, err = discordgo.New("Bot " + d.cfg.Token)
	if err != nil {
		return err
	}

	d.session.AddHandler(d.onMessage)
	d.session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	err = d.session.Open()
	if err != nil {
		utils.Logger.Error("Error opening Discord connection", zap.Error(err))
		return err
	}

	utils.Logger.Info("Discord Bot is now running")
	return nil
}

func (d *DiscordBot) Stop() {
	if d.session != nil {
		d.session.Close()
	}
}

func (d *DiscordBot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore self
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Filter guild if specified
	if d.cfg.GuildID != "" && m.GuildID != d.cfg.GuildID {
		return
	}

	isGroup := m.GuildID != ""

	d.handler(MessageEvent{
		Platform:   "discord",
		PlatformID: m.ChannelID,
		UserID:     m.Author.ID,
		Username:   m.Author.Username,
		Content:    m.Content,
		MsgType:    MsgTypeText,
		IsGroup:    isGroup,
	})
}

func (d *DiscordBot) SendMessage(targetID string, content string, isGroup bool) error {
	_, err := d.session.ChannelMessageSend(targetID, content)
	return err
}

func (d *DiscordBot) SendImage(targetID string, imageUrl string, isGroup bool) error {
	_, err := d.session.ChannelMessageSend(targetID, imageUrl)
	return err
}
