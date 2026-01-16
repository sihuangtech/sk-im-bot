package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	QQ       QQConfig       `mapstructure:"qq"`
	Discord  DiscordConfig  `mapstructure:"discord"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret         string `mapstructure:"secret"`
	ExpireDuration string `mapstructure:"expire_duration"`
}

type QQConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	WSURL       string `mapstructure:"ws_url"`
	AccessToken string `mapstructure:"access_token"`
}

type DiscordConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Token   string `mapstructure:"token"`
	GuildID string `mapstructure:"guild_id"`
}

type LLMConfig struct {
	Provider  string `mapstructure:"provider"`
	APIKey    string `mapstructure:"api_key"`
	BaseURL   string `mapstructure:"base_url"`
	Model     string `mapstructure:"model"`
	MaxTokens int    `mapstructure:"max_tokens"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
}

var GlobalConfig Config

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Automate env vars
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &GlobalConfig, nil
}
