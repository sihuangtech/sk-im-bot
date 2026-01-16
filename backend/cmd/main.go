package main

import (
	"fmt"

	"sk-im-bot/internal/api"
	"sk-im-bot/internal/bot"
	"sk-im-bot/internal/config"
	"sk-im-bot/internal/llm"
	"sk-im-bot/internal/model"
	"sk-im-bot/pkg/utils"
)

func main() {
	// 1. 加载配置文件
	// 默认从 "config/config.yaml" 读取
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("警告: 无法加载配置文件: %v。将尝试使用默认值或环境变量。\n", err)
		// 如果加载失败且没有默认配置，进行初始化
		if cfg == nil {
			cfg = &config.GlobalConfig
		}
	}

	// 2. 初始化核心工具类 (如日志记录器)
	// 根据配置中的日志级别进行初始化
	utils.InitLogger(cfg.Log.Level)
	defer utils.Logger.Sync() // 确保程序退出前刷新缓冲区

	// 3. 建立数据库连接
	// 初始化 PostgreSQL 数据库并执行自动迁移 (AutoMigrate)
	model.InitDB(cfg.Database)

	// 4. 启动 WebSocket 调度中心
	// 在独立协程中运行，负责管理前端管理界面的实时连接
	go api.WSHub.Run()

	// 5. 初始化大语言模型 (LLM) 客户端
	// 这里通常集成 OpenAI 服务，负责对话生成
	llmClient := llm.NewLLMClient(cfg.LLM)

	// 6. 初始化并启动机器人管理器
	// 管理器会启动已开启的平台（如 QQ 或 Discord）的机器人服务
	bot.InitManager(cfg, llmClient, api.BroadcastEvent)
	bot.Manager.Start()

	// 7. 配置并启动 Web API 服务器
	// 负责管理后台的 REST API 请求，如登录、统计信息获取等
	r := api.InitRouter()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	utils.Logger.Info(fmt.Sprintf("Web 服务器正在启动，监听地址: %s", addr))

	// 运行服务器，如果发生错误则记录致命日志并退出
	if err := r.Run(addr); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}
