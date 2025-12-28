package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
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
	ID          uint       `gorm:"primaryKey" json:"ID"`
	Name        string     `json:"Name"`
	CronExpr    string     `json:"CronExpr"`
	HostInfo    string     `gorm:"type:text" json:"HostInfo"` // MySQL/MariaDB 需要指定 type:text 以支持长文本
	Commands    string     `gorm:"type:text" json:"Commands"`
	Status      int        `json:"Status"` // 1: 启用, 0: 禁用
	RandomDelay int        `json:"RandomDelay"` // 随机延迟时间(分钟)
	LastRunTime *time.Time `json:"LastRunTime"`
	LastResult  string     `json:"LastResult"`
	ErrorLog    string     `gorm:"type:text" json:"ErrorLog"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
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
	dbType := os.Getenv("DB_TYPE") // sqlite (default) or mysql

	if dbType == "mysql" || dbType == "mariadb" {
		// MariaDB / MySQL 连接配置
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")

		if dbHost == "" { dbHost = "127.0.0.1" }
		if dbPort == "" { dbPort = "3306" }
		if dbUser == "" { dbUser = "root" }
		if dbName == "" { dbName = "webssh" }

		// DSN: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
			dbUser, dbPass, dbHost, dbPort, dbName)
		
		log.Printf("正在连接到 MariaDB/MySQL: %s:%s ...", dbHost, dbPort)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		// 默认 SQLite 逻辑
		if _, err := os.Stat("data"); os.IsNotExist(err) {
			err := os.Mkdir("data", 0755)
			if err != nil {
				log.Fatal("failed to create data directory: ", err)
			}
		}
		log.Println("正在使用 SQLite 数据库...")
		DB, err = gorm.Open(sqlite.Open("data/webssh.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// 自动迁移模式，自动创建或更新表结构
	err = DB.AutoMigrate(&User{}, &CronJob{}, &NotificationConfig{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	// 初始化默认管理员逻辑
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

	// 初始化空的通知配置
	var confCount int64
	DB.Model(&NotificationConfig{}).Count(&confCount)
	if confCount == 0 {
		DB.Create(&NotificationConfig{})
	}
}
