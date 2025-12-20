package model

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User 管理员用户
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"` // 存储明文或哈希
	CreatedAt time.Time `json:"created_at"`
}

// CronJob 定时任务
type CronJob struct {
	ID          uint   `gorm:"primaryKey" json:"ID"`
	Name        string `json:"Name"`
	CronExpr    string `json:"CronExpr"`
	HostInfo    string `json:"HostInfo"`
	Commands    string `json:"Commands"`
	Status      int    `json:"Status"` // 1: 启用, 0: 禁用
	RandomDelay int    `json:"RandomDelay"` // 随机延迟时间(分钟)
	LastRunTime *time.Time `json:"LastRunTime"`
	LastResult  string `json:"LastResult"`
	ErrorLog    string `json:"ErrorLog"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	ID          uint   `gorm:"primaryKey" json:"ID"`
	EnableEmail bool   `json:"enable_email"`
	EmailHost   string `json:"email_host"`
	EmailPort   int    `json:"email_port"`
	EmailUser   string `json:"email_user"`
	EmailPass   string `json:"email_pass"`
	EmailTo     string `json:"email_to"`
	
	EnableTg    bool   `json:"enable_tg"`
	TgBotToken  string `json:"tg_bot_token"`
	TgChatID    string `json:"tg_chat_id"`
	
	EnableBark  bool   `json:"enable_bark"`
	BarkUrl     string `json:"bark_url"`
}

// InitDB 初始化数据库连接及表结构
func InitDB() {
	var err error
	
	// 确保数据目录存在
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err := os.Mkdir("data", 0755)
		if err != nil {
			log.Fatal("failed to create data directory: ", err)
		}
	}

	// 连接 SQLite 数据库，文件名为 data/webssh.db
	DB, err = gorm.Open(sqlite.Open("data/webssh.db"), &gorm.Config{})
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
