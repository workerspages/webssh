package model

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User 管理员用户
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique"`
	Password  string    `json:"-"` // 存储明文或哈希，本示例为演示方便存明文，生产环境请Bcrypt
	CreatedAt time.Time
}

// CronJob 定时任务
type CronJob struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	CronExpr    string // Cron表达式
	HostInfo    string // Base64编码的SSH连接信息
	Commands    string // JSON字符串 ["cmd1", "cmd2"]
	Status      bool   // true: 启用, false: 禁用
	LastRunTime *time.Time
	LastResult  string // 成功/失败
	ErrorLog    string // 具体的错误日志
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	ID          uint `gorm:"primaryKey"`
	EmailHost   string
	EmailPort   int
	EmailUser   string
	EmailPass   string
	EmailTo     string
	TgBotToken  string
	TgChatID    string
	EnableEmail bool
	EnableTg    bool
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("webssh.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// 自动迁移模式
	DB.AutoMigrate(&User{}, &CronJob{}, &NotificationConfig{})

	// 初始化默认管理员 (admin/admin123)
	var count int64
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		DB.Create(&User{Username: "admin", Password: "admin123"})
		log.Println("Default user created: admin / admin123")
	}
	
	// 初始化空的配置
	var confCount int64
	DB.Model(&NotificationConfig{}).Count(&confCount)
	if confCount == 0 {
		DB.Create(&NotificationConfig{})
	}
}
