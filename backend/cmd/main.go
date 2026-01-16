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
	// 1. Load Config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v. Using defaults/env.\n", err)
		// We proceed, as config might be empty initially
		if cfg == nil {
			cfg = &config.GlobalConfig
		}
	}

	// 2. Init Utils (Logger)
	utils.InitLogger(cfg.Log.Level)
	defer utils.Logger.Sync()

	// 3. Init Database
	model.InitDB(cfg.Database)

	// 4. Init WebSocket Hub
	go api.WSHub.Run()

	// 5. Init LLM
	llmClient := llm.NewLLMClient(cfg.LLM)

	// 6. Init Bot Manager
	bot.InitManager(cfg, llmClient, api.BroadcastEvent)
	bot.Manager.Start()

	// 7. Start Web Server
	r := api.InitRouter()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	utils.Logger.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := r.Run(addr); err != nil {
		utils.Logger.Fatal(err.Error())
	}
}
