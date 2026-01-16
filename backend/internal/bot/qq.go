package bot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"sk-im-bot/internal/config"
	"sk-im-bot/pkg/utils"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type QQBot struct {
	cfg         config.QQConfig
	conn        *websocket.Conn
	mu          sync.Mutex
	handler     func(MessageEvent)
	isConnected bool
}

func NewQQBot(cfg config.QQConfig, handler func(MessageEvent)) *QQBot {
	return &QQBot{
		cfg:     cfg,
		handler: handler,
	}
}

func (q *QQBot) Start() error {
	utils.Logger.Info("Connecting to QQ OneBot...", zap.String("url", q.cfg.WSURL))

	// Basic retry loop
	for {
		err := q.connect()
		if err != nil {
			utils.Logger.Error("QQ Connect failed, retrying in 5s", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}

		// Reading loop
		for {
			_, message, err := q.conn.ReadMessage()
			if err != nil {
				utils.Logger.Error("QQ Read error", zap.Error(err))
				q.isConnected = false
				break
			}
			q.parseMessage(message)
		}
	}
}

func (q *QQBot) connect() error {
	u, err := url.Parse(q.cfg.WSURL)
	if err != nil {
		return err
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	q.conn = c
	q.isConnected = true
	utils.Logger.Info("Connected to QQ OneBot")
	return nil
}

func (q *QQBot) Stop() {
	if q.conn != nil {
		q.conn.Close()
	}
}

// OneBot 11 simple payload
type OneBotEvent struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"` // private, group
	SubType     string `json:"sub_type"`
	UserId      int64  `json:"user_id"`
	GroupId     int64  `json:"group_id"`
	RawMessage  string `json:"raw_message"`
	Sender      struct {
		Nickname string `json:"nickname"`
	} `json:"sender"`
}

func (q *QQBot) parseMessage(data []byte) {
	var event OneBotEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return // ignore non-json or echo
	}

	if event.PostType == "message" {
		isGroup := event.MessageType == "group"
		targetID := fmt.Sprintf("%d", event.UserId)
		if isGroup {
			targetID = fmt.Sprintf("%d", event.GroupId)
		}

		q.handler(MessageEvent{
			Platform:   "qq",
			PlatformID: targetID,
			UserID:     fmt.Sprintf("%d", event.UserId),
			Username:   event.Sender.Nickname,
			Content:    event.RawMessage,
			MsgType:    MsgTypeText,
			IsGroup:    isGroup,
		})
	}
}

type OneBotAction struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Echo   string      `json:"echo"`
}

func (q *QQBot) SendMessage(targetID string, content string, isGroup bool) error {
	if !q.isConnected {
		return fmt.Errorf("QQ not connected")
	}

	action := "send_private_msg"
	params := map[string]interface{}{
		"user_id": targetID,
		"message": content,
	}

	if isGroup {
		action = "send_group_msg"
		params = map[string]interface{}{
			"group_id": targetID,
			"message":  content,
		}
	}

	payload := OneBotAction{
		Action: action,
		Params: params,
		Echo:   fmt.Sprintf("send_%d", time.Now().UnixNano()),
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	return q.conn.WriteJSON(payload)
}

func (q *QQBot) SendImage(targetID string, imageUrl string, isGroup bool) error {
	// Construct CQ code [CQ:image,file=http://...]
	cqCode := fmt.Sprintf("[CQ:image,file=%s]", imageUrl)
	return q.SendMessage(targetID, cqCode, isGroup)
}
