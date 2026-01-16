package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 高性能日志记录组件单例，系统统一使用 Zap 日志库输出
var Logger *zap.Logger

// InitLogger 初始化强类型的、高性能的结构化日志组件
func InitLogger(level string) {
	// 创建生产级别的默认配置
	config := zap.NewProductionConfig()

	// 从字符串动态配置日志的拦截级别 (debug -> 打印所有调试信息)
	if level == "debug" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		// 普通运行状态仅输出关键业务流程及报错
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 时间戳格式化为符合中国人阅读习惯的 ISO 标准
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	// 构建最终的 Logger 实例
	Logger, err = config.Build()
	if err != nil {
		// 日志驱动故障属于基础性错误，直接 panic 暴露
		panic("日志系统无法启动, 请检查配置环境")
	}
}

// Claims 自定义携带的令牌载荷，通过此结构提取当前请求对应的后台操作者信息
type Claims struct {
	UserID               uint   `json:"user_id"` // 当前鉴权成功的管理人员 UID
	Role                 string `json:"role"`    // 该人员被分配的权限角色
	jwt.RegisteredClaims        // 混入 JWT 标准规定的预定义载荷 (exp, iat 等)
}

// GenerateToken 执行 JWT 签发逻辑。使用 HS256 对称加密算法构建安全令牌
func GenerateToken(userID uint, role string, secret string, duration time.Duration) (string, error) {
	// 构造令牌携带的详细内容
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置令牌自动过期时刻
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			// 设置令牌签发时刻
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// 选用 HMAC-SHA256 方式进行摘要签名，并注入配置中的 Secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret)) // 输出加密后的最终字符串
}

// ParseToken 实现令牌反解与校验。由于使用了固定密钥，这里会验证签名完整性
func ParseToken(tokenString string, secret string) (*Claims, error) {
	// 将原始字符串解码，并尝试反序列化为我们自定义的 Claims 模型
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 回调函数提供校验所需的密钥
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证：1. 模型解析成功 2. 内部标准载荷校验通过（如未过期）
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
