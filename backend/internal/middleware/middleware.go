package middleware

import (
	"net/http"
	"strings"

	"sk-im-bot/internal/config"
	"sk-im-bot/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 强力且安全的身份验证中间件，拦截所有受保护路由
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 HTTP 头部读取 Authorization 段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 未发现令牌请求头，直接拦截
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请在请求头中携带身份令牌 (Authorization)"})
			c.Abort() // 终止当前请求链路后续处理
			return
		}

		// 验证令牌前缀格式是否符合 RFC 6750 规范 (Bearer Tokens)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权格式不合法，请使用 'Bearer <token>'"})
			c.Abort()
			return
		}

		// 解析并验证 JWT 令牌的签名和过期时间
		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString, config.GlobalConfig.JWT.Secret)
		if err != nil {
			// 令牌失效或篡改
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已失效或非法，请重新登录"})
			c.Abort()
			return
		}

		// 将解析后的关键身份元数据 (UID, Role) 注入上下文，供后续 Handler 获取
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next() // 校验通过，执行后续逻辑
	}
}

// CORSMiddleware 解决跨域限制策略的中间件，由于本项目前后端通常异端口，必须开启
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许跨越的源地址（* 代表全部匹配）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 支持复杂请求携带凭证（如 Cookie 等）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// 定义前端代码可以自由使用的 HTTP 响应头和请求头列表
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 允许的交互请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 所有的复杂请求 (例如带 JSON 或 Header 的 POST) 都会先发一个 OPTIONS 请求作为嗅探
		if c.Request.Method == "OPTIONS" {
			// 直接对此嗅探请求返回 204 无损状态码
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
