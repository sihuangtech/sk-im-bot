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

// QQBot 适配器，通过 OneBot 11 协议对接 QQ 客户端 (如 go-cqhttp)
type QQBot struct {
	cfg         config.QQConfig    // QQ 配置（服务器地址、访问链等）
	conn        *websocket.Conn    // 后端与 OneBot 端的 WebSocket 连接
	mu          sync.Mutex         // 互斥锁，确保并发写操作安全
	handler     func(MessageEvent) // 消息接收回调处理逻辑
	isConnected bool               // 运行时的连接状态标记
}

// NewQQBot 构造一个全新的 QQ 机器人适配器
func NewQQBot(cfg config.QQConfig, handler func(MessageEvent)) *QQBot {
	return &QQBot{
		cfg:     cfg,
		handler: handler,
	}
}

// Start 启动对 OneBot 服务器的长连接监听，包含掉线重连逻辑
func (q *QQBot) Start() error {
	utils.Logger.Info("正在尝试连接到 QQ OneBot 服务...", zap.String("url", q.cfg.WSURL))

	// 无限重连循环，确保服务高可用
	for {
		err := q.connect()
		if err != nil {
			utils.Logger.Error("QQ 连接失败，5秒后重试...", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}

		// 读取循环
		for {
			_, message, err := q.conn.ReadMessage()
			if err != nil {
				utils.Logger.Error("QQ 连接断开 (读取错误)", zap.Error(err))
				q.isConnected = false
				break
			}
			// 异步处理收到的原始二进制/Json 数据
			go q.parseMessage(message)
		}
	}
}

// connect 实际执行向 OneBot 协议端握手
func (q *QQBot) connect() error {
	u, err := url.Parse(q.cfg.WSURL)
	if err != nil {
		return err
	}

	// 使用 gorilla/websocket 向目标地址（通常是端点侧的 WebSocket 正向监听）拨号
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	q.conn = c
	q.isConnected = true
	utils.Logger.Info("成功建立 QQ OneBot WebSocket 通讯")
	return nil
}

// Stop 优雅关闭机器人连接
func (q *QQBot) Stop() {
	if q.conn != nil {
		q.conn.Close()
	}
}

// OneBotEvent 映射 OneBot 11 上报的 JSON 原始字段
type OneBotEvent struct {
	PostType    string `json:"post_type"`    // 事件类型 (message, notice, request, meta_event)
	MessageType string `json:"message_type"` // 消息子类型 (private, group)
	SubType     string `json:"sub_type"`
	UserId      int64  `json:"user_id"`     // 发送者 QQ 号
	GroupId     int64  `json:"group_id"`    // 群组号
	RawMessage  string `json:"raw_message"` // 原始文本消息，含 CQ 码
	Sender      struct {
		Nickname string `json:"nickname"` // 发送者昵称
	} `json:"sender"`
}

// parseMessage 将接收到的 OneBot JSON 数据解析为内部统一结构
func (q *QQBot) parseMessage(data []byte) {
	var event OneBotEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return // 忽略控制包、Meta 数据或非格式化报文
	}

	// 仅监听消息上报
	if event.PostType == "message" {
		isGroup := event.MessageType == "group"
		// 目标 ID：群聊取群号，私聊取对方 QQ 号
		targetID := fmt.Sprintf("%d", event.UserId)
		if isGroup {
			targetID = fmt.Sprintf("%d", event.GroupId)
		}

		// 调用上层管理逻辑的回调函数
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

// OneBotAction 定义发送给端点的 API 调用结构
type OneBotAction struct {
	Action string      `json:"action"` // 动作名称
	Params interface{} `json:"params"` // 参数负载
	Echo   string      `json:"echo"`   // 回显标识
}

// SendMessage 向指定的目标发送 QQ 消息 (支持私聊和群聊)
func (q *QQBot) SendMessage(targetID string, content string, isGroup bool) error {
	if !q.isConnected {
		return fmt.Errorf("QQ 机器人当前未在线")
	}

	action := "send_private_msg"
	params := map[string]interface{}{
		"user_id": targetID,
		"message": content,
	}

	// 如果是群聊则切换 API 动作为发送群消息
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
		Echo:   fmt.Sprintf("send_%d", time.Now().UnixNano()), // 使用纳秒级时间确保 Echo 唯一
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	// 将 API 调用请求以 JSON 形式通过长连接发出
	return q.conn.WriteJSON(payload)
}

// SendImage 利用 CQ 码发送富媒体图片消息
func (q *QQBot) SendImage(targetID string, imageUrl string, isGroup bool) error {
	// 构建 OneBot 规范的图片 CQ 码负载
	cqCode := fmt.Sprintf("[CQ:image,file=%s]", imageUrl)
	return q.SendMessage(targetID, cqCode, isGroup)
}
