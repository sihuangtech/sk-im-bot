package api

import (
	"sk-im-bot/internal/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 路由器，设置全局中间件和 API 路由
func InitRouter() *gin.Engine {
	// 创建一个带默认中间件（Logger 和 Recovery）的 Gin 引擎
	r := gin.Default()

	// 挂载跨域 (CORS) 中间件，允许前端跨域访问 API
	r.Use(middleware.CORSMiddleware())

	// ---- 公开路由 (无需鉴权) ----

	// 管理员登录接口
	r.POST("/api/login", Login)

	// WebSocket 实时监控连接接口
	r.GET("/ws", WSHandler)

	// ---- 受保护路由 (需要 JWT 鉴权) ----

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware()) // 应用身份验证中间件
	{
		// 获取/更新系统配置项
		api.GET("/config", GetConfig)
		api.POST("/config", UpdateConfig)

		// 获取消息历史及会话管理数据
		api.GET("/messages", GetMessages)
		api.GET("/sessions", GetSessions)
	}

	return r
}
