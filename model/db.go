package model

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User 管理员用户
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique"`
	Password  string    `json:"-"` // 存储明文或哈希，本示例为演示方便存明文，生产环境建议使用Bcrypt加密
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

// InitDB 初始化数据库连接及表结构
func InitDB() {
	var err error
	// 连接 SQLite 数据库，文件名为 webssh.db
	DB, err = gorm.Open(sqlite.Open("webssh.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// 自动迁移模式，自动创建或更新表结构
	err = DB.AutoMigrate(&User{}, &CronJob{}, &NotificationConfig{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	// 初始化默认管理员逻辑
	// 优先读取环境变量 USER 和 PASS，如果没有则使用默认值 admin/admin123
	var count int64
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		defaultUser := os.Getenv("USER")
		if defaultUser == "" {
			defaultUser = "admin"
		}

		defaultPass := os.Getenv("PASS")
		if defaultPass == "" {
			defaultPass = "admin123"
		}

		user := User{Username: defaultUser, Password: defaultPass}
		if result := DB.Create(&user); result.Error != nil {
			log.Printf("创建初始管理员失败: %v\n", result.Error)
		} else {
			log.Printf("初始化管理员账号成功: %s / %s (请注意保护账号安全)\n", defaultUser, defaultPass)
		}
	}

	// 初始化空的通知配置（保证数据库中至少有一条配置记录，ID通常为1）
	var confCount int64
	DB.Model(&NotificationConfig{}).Count(&confCount)
	if confCount == 0 {
		DB.Create(&NotificationConfig{})
	}
}
