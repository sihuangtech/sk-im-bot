package model

import (
	"time"

	"gorm.io/gorm"
)

// User 定义了管理系统的用户信息
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`        // 主键ID
	Username  string         `gorm:"uniqueIndex" json:"username"` // 唯一用户名
	Password  string         `json:"-"`                           // 密码哈希（不返回给前端）
	Role      string         `json:"role"`                        // 角色: admin, user
	CreatedAt time.Time      `json:"created_at"`                  // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`                  // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`              // 软删除支持
}

// Session 代表机器人与用户或群组的一个会话实例
type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`     // 会话内部ID
	Platform     string    `gorm:"index" json:"platform"`    // 平台类型: qq, discord
	PlatformID   string    `gorm:"index" json:"platform_id"` // 平台侧的ID (群ID或用户ID)
	PlatformName string    `json:"platform_name"`            // 平台侧显示的名称
	LastActive   time.Time `json:"last_active"`              // 最后活跃时间
}

// Message 存储所有的聊天历史记录
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`    // 消息ID
	SessionID uint      `gorm:"index" json:"session_id"` // 所属会话ID
	Sender    string    `json:"sender"`                  // 发送者名称 (user 或 bot)
	Content   string    `json:"content"`                 // 消息内容文本
	MsgType   string    `json:"msg_type"`                // 消息类型: text, image
	RawData   string    `json:"raw_data"`                // 原始JSON数据备份
	CreatedAt time.Time `json:"created_at"`              // 接收/发送时间
}

// Config 存储系统动态配置（数据库持久化版本）
type Config struct {
	Key         string `gorm:"primaryKey" json:"key"` // 配置键
	Value       string `json:"value"`                 // 配置值
	Description string `json:"description"`           // 描述信息
}

// Blacklist 存储黑名单用户信息
type Blacklist struct {
	ID        uint      `gorm:"primaryKey" json:"id"` // 主键
	Platform  string    `json:"platform"`             // 限定平台
	TargetID  string    `json:"target_id"`            // 目标用户识别码
	Reason    string    `json:"reason"`               // 拉黑原因
	CreatedAt time.Time `json:"created_at"`           // 拉黑时间
}
