package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"webssh/model"
	"webssh/service"

	"github.com/gin-gonic/gin"
)

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	var form struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&form); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	token, err := service.Login(form.Username, form.Password)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "success", "token": token})
}

// GetCronJobs 获取任务列表
func GetCronJobs(c *gin.Context) {
	var jobs []model.CronJob
	model.DB.Order("id desc").Find(&jobs)
	c.JSON(200, gin.H{"code": 200, "data": jobs})
}

// AddCronJob 添加/更新任务
func AddCronJob(c *gin.Context) {
	var job model.CronJob
	if err := c.BindJSON(&job); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	
	// 验证 SSH Info 是否有效 (略，可复用 core.CheckSSH)
	
	if job.ID > 0 {
		model.DB.Save(&job)
	} else {
		model.DB.Create(&job)
	}
	service.ReloadJobs() // 重载调度器
	c.JSON(200, gin.H{"code": 200, "msg": "saved"})
}

// DeleteCronJob 删除任务
func DeleteCronJob(c *gin.Context) {
	id := c.Param("id")
	model.DB.Delete(&model.CronJob{}, id)
	service.ReloadJobs()
	c.JSON(200, gin.H{"code": 200, "msg": "deleted"})
}

// RunCronJobManually 手动触发
func RunCronJobManually(c *gin.Context) {
	id := c.Param("id")
	var job model.CronJob
	if err := model.DB.First(&job, id).Error; err != nil {
		c.JSON(400, gin.H{"msg": "job not found"})
		return
	}
	go service.RunJob(&job)
	c.JSON(200, gin.H{"code": 200, "msg": "started"})
}

// GetNotifyConfig 获取配置
func GetNotifyConfig(c *gin.Context) {
	var conf model.NotificationConfig
	model.DB.First(&conf)
	c.JSON(200, gin.H{"code": 200, "data": conf})
}

// SaveNotifyConfig 保存配置
func SaveNotifyConfig(c *gin.Context) {
	var conf model.NotificationConfig
	if err := c.BindJSON(&conf); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	// 总是更新ID为1的记录
	conf.ID = 1
	model.DB.Save(&conf)
	c.JSON(200, gin.H{"code": 200, "msg": "saved"})
}

// TestNotify 测试通知
func TestNotify(c *gin.Context) {
	service.SendNotification("WebSSH 测试通知", "这是一条测试消息，如果您收到它，说明通知配置正确。")
	c.JSON(200, gin.H{"code": 200, "msg": "sent"})
}
