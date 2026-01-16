package model

import (
	"fmt"
	"log"

	"sk-im-bot/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 全局数据库 ORM 操作实例，本项目各模块通过此句柄访问数据库
var DB *gorm.DB

// InitDB 根据用户提供的配置初始化 PostgreSQL 的长连接池
func InitDB(cfg config.DatabaseConfig) {
	// 拼接 PostgreSQL 标准 DSN 连接格式
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	var err error
	// 使用 GORM 驱动开启连接池，并注入 PostgreSQL 专用驱动
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果数据库连不上，属于致命错误，强制退出程序
		log.Fatalf("无法连接到 PostgreSQL 数据库集群: %v", err)
	}

	// 自动迁移 (AutoMigrate)
	// 本功能会自动对比代码结构体与表结构的差异，执行建表或字段新增操作
	// 注意：生产环境下严禁使用该功能删除数据列
	err = DB.AutoMigrate(&User{}, &Session{}, &Message{}, &Config{}, &Blacklist{})
	if err != nil {
		log.Fatalf("执行数据库模型自动迁移失败: %v", err)
	}

	log.Println("数据库层初始化完成，所有数据模型同步完毕")
}
