package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"webssh/core"
	"webssh/model"

	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

func StartCron() {
	Cron = cron.New()
	Cron.Start()
	ReloadJobs()
}

// ReloadJobs 从数据库重新加载所有任务
func ReloadJobs() {
	// 清除旧任务（简单粗暴方式，重启所有）
	items := Cron.Entries()
	for _, entry := range items {
		Cron.Remove(entry.ID)
	}

	var jobs []model.CronJob
	model.DB.Where("status = ?", 1).Find(&jobs)

	for _, job := range jobs {
		j := job // copy for closure
		_, err := Cron.AddFunc(j.CronExpr, func() {
			RunJob(&j)
		})
		if err != nil {
			log.Printf("Failed to add cron job %s: %v", j.Name, err)
		}
	}
}

// RunJob 立即执行任务
func RunJob(job *model.CronJob) {
	log.Printf("Starting job: %s", job.Name)
	
	// 更新开始状态
	now := time.Now()
	
	// 1. 解析 SSH 信息
	client, err := core.DecodedMsgToSSHClient(job.HostInfo)
	if err != nil {
		handleJobResult(job, now, "失败", fmt.Sprintf("SSH配置解析失败: %v", err))
		return
	}

	// 2. 解析命令列表
	var cmds []string
	if err := json.Unmarshal([]byte(job.Commands), &cmds); err != nil {
		handleJobResult(job, now, "失败", "命令格式解析失败")
		return
	}

	// 3. 执行
	output, err := client.RunBatchTasks(cmds)
	
	// 4. 处理结果
	if err != nil {
		handleJobResult(job, now, "失败", output)
		// 发送通知
		SendNotification(
			fmt.Sprintf("【WebSSH 报警】任务执行失败: %s", job.Name),
			fmt.Sprintf("服务器: %s\n执行时间: %s\n\n错误日志:\n%s", client.Hostname, now.Format(time.RFC3339), output),
		)
	} else {
		handleJobResult(job, now, "成功", output)
	}
}

func handleJobResult(job *model.CronJob, runTime time.Time, result, logStr string) {
	// 重新从DB获取以避免并发覆盖
	var currentJob model.CronJob
	model.DB.First(&currentJob, job.ID)
	currentJob.LastRunTime = &runTime
	currentJob.LastResult = result
	currentJob.ErrorLog = logStr
	model.DB.Save(&currentJob)
}
