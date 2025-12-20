package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户表
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string // 实际生产请存储 Hash
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	gorm.Model
	Type      string // "email" or "telegram"
	ConfigJson string // 存储JSON格式的配置 (SMTP host, port... / TG Token, ChatID)
	IsEnabled bool
}

// CronJob 定时任务
type CronJob struct {
	gorm.Model
	Name        string
	CronExpr    string // e.g., "0 0 * * *"
	HostInfo    string // 存储 SSH 连接信息的 JSON (IP, Port, User, Pass/Key)
	Commands    string // 存储命令列表 JSON ["cd /var", "rm -rf tmp"]
	Status      int    // 0: 停止, 1: 运行中
	LastRunTime *time.Time
	LastResult  string
}

// 初始化数据库
func InitDB(dbPath string) (*gorm.DB, error) {
    // 这里具体实现连接 sqlite 代码
    // db.AutoMigrate(&User{}, &NotificationConfig{}, &CronJob{})
    return nil, nil // 占位
}
