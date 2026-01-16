package api

import (
	"net/http"
	"time"

	"sk-im-bot/internal/config"
	"sk-im-bot/internal/model"
	"sk-im-bot/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Login 处理后台管理员登录请求
func Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// 绑定并校验 JSON 输入
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据"})
		return
	}

	// 演示目的：使用硬编码的凭据 (admin/admin)。
	// 在正式生产环境中，应查库比对哈希后的密码。
	if body.Username == "admin" && body.Password == "admin" {
		// 校验成功，生成 JWT Token，有效期 24 小时
		token, err := utils.GenerateToken(1, "admin", config.GlobalConfig.JWT.Secret, 24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token 生成失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	// 校验失败返回未授权状态
	c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
}

// GetConfig 获取当前系统的全局配置项
func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.GlobalConfig)
}

// UpdateConfig 在线更新系统配置
func UpdateConfig(c *gin.Context) {
	var newConfig config.Config
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 在内存中实时更新配置，正式应用建议同步写入文件或持久化到 DB
	config.GlobalConfig = newConfig
	c.JSON(http.StatusOK, gin.H{"status": "配置已更新"})
}

// GetMessages 获取最近的历史消息记录
func GetMessages(c *gin.Context) {
	var messages []model.Message
	// 从数据库查询最近的 50 条消息，按时间倒序排列
	model.DB.Order("created_at desc").Limit(50).Find(&messages)
	c.JSON(http.StatusOK, messages)
}

// GetSessions 获取系统中参与对话的所有活跃会话列表
func GetSessions(c *gin.Context) {
	var sessions []model.Session
	// 查询全量会话信息
	model.DB.Find(&sessions)
	c.JSON(http.StatusOK, sessions)
}
