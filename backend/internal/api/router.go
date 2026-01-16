package api

import (
	"sk-im-bot/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.POST("/api/login", Login)
	r.GET("/ws", WSHandler)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/config", GetConfig)
		api.POST("/config", UpdateConfig)
		api.GET("/messages", GetMessages)
		api.GET("/sessions", GetSessions)
	}

	return r
}
