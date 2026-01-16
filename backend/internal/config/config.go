package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置根结构体，映射 YAML 配置文件中的全部树状字段
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	QQ       QQConfig       `mapstructure:"qq"`
	Discord  DiscordConfig  `mapstructure:"discord"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Log      LogConfig      `mapstructure:"log"`
	Admin    AdminConfig    `mapstructure:"admin"`

	// Runtime only, loaded from llm_providers.yaml
	LLMProviders map[string]LLMProviderConfig `mapstructure:"-"`
}

// LLMProviderConfig 定义单个 LLM 提供商的连接预设
type LLMProviderConfig struct {
	BaseURL      string   `mapstructure:"base_url"`
	DefaultModel string   `mapstructure:"default_model"`
	Models       []string `mapstructure:"models"`
}

// AdminConfig 定义后台管理账号配置
type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ServerConfig 定义管理系统的 Web 服务选项
type ServerConfig struct {
	Port int    `mapstructure:"port"` // 监听端口 (默认 8888)
	Mode string `mapstructure:"mode"` // 运行模式 (debug 或 release)
}

// DatabaseConfig 定义 PostgreSQL 连接凭证及地址
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig 访问令牌验证相关参数配置
type JWTConfig struct {
	Secret         string `mapstructure:"secret"`          // 加密签名密钥
	ExpireDuration string `mapstructure:"expire_duration"` // 授权过期周期 (例如 24h)
}

// QQConfig QQ 接入（通常针对 go-cqhttp）的具体设置
type QQConfig struct {
	Enabled     bool   `mapstructure:"enabled"`      // 是否激活此模块
	WSURL       string `mapstructure:"ws_url"`       // WebSocket 长连地址 (e.g. ws://localhost:8080)
	AccessToken string `mapstructure:"access_token"` // OneBot 安全访问凭据 (如有)
}

// DiscordConfig Discord 服务接入参数
type DiscordConfig struct {
	Enabled bool   `mapstructure:"enabled"`  // 是否激活此模块
	Token   string `mapstructure:"token"`    // 机器人应用 Token (Bot Token)
	GuildID string `mapstructure:"guild_id"` // 限制监听的特定服务器 ID (选填)
}

// LLMConfig 大语言模型（如 OpenAI）调用的鉴权与参数集
type LLMConfig struct {
	Provider  string `mapstructure:"provider"`   // 厂商标识 (openai/claude 等)
	APIKey    string `mapstructure:"api_key"`    // API 访问密钥
	BaseURL   string `mapstructure:"base_url"`   // 访问网关 (支持中转代理由此输入)
	Model     string `mapstructure:"model"`      // 指定模型版本 (e.g. gpt-4)
	MaxTokens int    `mapstructure:"max_tokens"` // 限制单次回复的最大 Token 数
}

// LogConfig 系统运行日志存储配置
type LogConfig struct {
	Level    string `mapstructure:"level"`    // 记录等级 (info, error, debug)
	Filename string `mapstructure:"filename"` // 导出文件名
}

// GlobalConfig 内存中持有的实时配置快照，系统各模块共享读取
var GlobalConfig Config

// LoadConfig 使用 Viper 库加载磁盘上的 YAML 配置文件并监听环境映射
func LoadConfig(path string) (*Config, error) {
	// 1. 设置环境变量替换规则: 将配置中的 "." 替换为 "_"
	// 例如: Server.Port -> SERVER_PORT
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 2. 自动加载系统环境变量
	viper.AutomaticEnv()

	// 3. 尝试加载 .env 文件 (仅用于本地开发，生产环境通常直接注入环境变量)
	// 如果指定路径为空，默认尝试从当前目录加载 .env
	if path == "" {
		viper.AddConfigPath(".")
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
	} else {
		// 如果指定了具体文件（如 config/config.yaml），则加载它
		// 但为了支持 env 覆盖，我们仍保留读取 .env 的能力
		// 此处逻辑修改为：优先读取 .env 文件（如果存在），再整合 System Env
		// 原 yaml 读取逻辑保留作为 fallback 或基础配置，用户可留空
		viper.SetConfigFile(path)
	}

	// 尝试读取配置 (文件)
	if err := viper.ReadInConfig(); err != nil {
		// 如果是“未找到配置文件”错误，且我们主要依赖环境变量，则不应 panic/error
		// 但 Viper 对于 SetConfigFile 如果文件不存在会报错
		// 我们改为: 记录警告但不中断，除非完全没有配置源
		fmt.Printf("提示: 未在 %s 找到配置文件或 .env，将完全依赖系统环境变量: %v\n", path, err)
	}

	// 将解析结果解包到内存结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return nil, fmt.Errorf("反序列化配置失败: %w", err)
	}

	// 4. 单独加载 LLM 模型提供商配置 (backend/config/llm_providers.yaml)
	// 我们使用一个新的 viper 实例来避免与主配置混淆，或者将其合并到 map 中
	providerViper := viper.New()
	providerViper.SetConfigFile("config/llm_providers.yaml")
	providerViper.SetConfigType("yaml")

	if err := providerViper.ReadInConfig(); err == nil {
		var providerConfig struct {
			Providers map[string]LLMProviderConfig `mapstructure:"providers"`
		}
		if err := providerViper.Unmarshal(&providerConfig); err == nil {
			GlobalConfig.LLMProviders = providerConfig.Providers
			fmt.Println("成功加载 LLM 提供商预设配置")
		} else {
			fmt.Printf("警告: 解析 llm_providers.yaml 失败: %v\n", err)
		}
	} else {
		// 可能是生产环境路径不同，尝试从当前目录加载
		providerViper.SetConfigFile("llm_providers.yaml")
		if err := providerViper.ReadInConfig(); err == nil {
			var providerConfig struct {
				Providers map[string]LLMProviderConfig `mapstructure:"providers"`
			}
			if err := providerViper.Unmarshal(&providerConfig); err == nil {
				GlobalConfig.LLMProviders = providerConfig.Providers
			}
		}
	}

	return &GlobalConfig, nil
}
