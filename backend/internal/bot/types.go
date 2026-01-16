package bot

import (
	"math/rand"
	"time"
)

type MsgType string

const (
	MsgTypeText  MsgType = "text"
	MsgTypeImage MsgType = "image"
)

type MessageEvent struct {
	Platform     string
	PlatformID   string // ChannelID for Discord, GroupID/UserID for QQ
	UserID       string
	Username     string
	Content      string
	MsgType      MsgType
	IsGroup      bool
}

type BotAdapter interface {
	Start() error
	Stop()
	SendMessage(targetID string, content string, isGroup bool) error
	SendImage(targetID string, imageUrl string, isGroup bool) error
}

// RandomDelay sleeps for a random duration between min and max milliseconds
func RandomDelay(minMs, maxMs int) {
	if minMs >= maxMs {
		return
	}
	duration := time.Duration(rand.Intn(maxMs-minMs)+minMs) * time.Millisecond
	time.Sleep(duration)
}
