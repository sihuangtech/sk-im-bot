package api

import (
	"net/http"
	"time"

	"sk-im-bot/internal/config"
	"sk-im-bot/internal/model"
	"sk-im-bot/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// For demo purpose, hardcoded check. Real world use DB check/Hash
	if body.Username == "admin" && body.Password == "admin" {
		token, err := utils.GenerateToken(1, "admin", config.GlobalConfig.JWT.Secret, 24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.GlobalConfig)
}

func UpdateConfig(c *gin.Context) {
	var newConfig config.Config
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// In real app, save to file or DB and reload.
	config.GlobalConfig = newConfig
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func GetMessages(c *gin.Context) {
	var messages []model.Message
	// Pagination logic here
	model.DB.Order("created_at desc").Limit(50).Find(&messages)
	c.JSON(http.StatusOK, messages)
}

func GetSessions(c *gin.Context) {
	// Simple unique session aggregation or separate table
	var sessions []model.Session
	model.DB.Find(&sessions)
	c.JSON(http.StatusOK, sessions)
}
