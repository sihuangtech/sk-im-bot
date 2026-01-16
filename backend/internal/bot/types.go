package bot

import (
	"math/rand"
	"time"
)

// MsgType 定义支持的消息媒体类型
type MsgType string

const (
	MsgTypeText  MsgType = "text"  // 普通文本
	MsgTypeImage MsgType = "image" // 图片数据
)

// MessageEvent 定义了所有机器人适配器 (QQ/Discord) 的通用消息载荷格式
// 后端通过此结构体抹平不同平台报文的差异
type MessageEvent struct {
	Platform   string  // 平台标识：qq, discord
	PlatformID string  // 平台目标 ID (群号、频道 ID、或用户识别码)
	UserID     string  // 消息发送方的唯一 ID
	Username   string  // 发送方显示的屏幕昵称
	Content    string  // 消息文本正文
	MsgType    MsgType // 消息类型
	IsGroup    bool    // 是否属于群组/大群环境
}

// BotAdapter 平台适配器接口定义。新对接平台（如 Telegram 或微信）必须实现这些方法
type BotAdapter interface {
	Start() error                                                    // 启动监听任务
	Stop()                                                           // 安全停止进程
	SendMessage(targetID string, content string, isGroup bool) error // 执行文本回复发送
	SendImage(targetID string, imageUrl string, isGroup bool) error  // 执行图片附件发送
}

// RandomDelay 用于模拟真人行为。在指定范围内产生随机毫秒级的阻塞延迟
// 参数 minMs, maxMs 为毫秒单位的下限和上限
func RandomDelay(minMs, maxMs int) {
	if minMs >= maxMs {
		return
	}
	// 计算随机偏移量
	duration := time.Duration(rand.Intn(maxMs-minMs)+minMs) * time.Millisecond
	time.Sleep(duration) // 阻塞当前协程指定时长
}
