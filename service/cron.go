package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"webssh/core"
	"webssh/model"

	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

func StartCron() {
	Cron = cron.New()
	Cron.Start()
	log.Printf("Cron Manager Started. Current Server Time: %s", time.Now().Format(time.RFC3339))
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
		id, err := Cron.AddFunc(j.CronExpr, func() {
			RunJob(&j)
		})
		if err != nil {
			log.Printf("Failed to add cron job %s: %v", j.Name, err)
		} else {
			log.Printf("Registered job [%s] with schedule: %s (ID: %d)", j.Name, j.CronExpr, id)
		}
	}
}

// SSHInfo represents the structure of host information
type SSHInfo struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	Passphrase string `json:"passphrase"`
	LoginType  int    `json:"loginType"` // 0 for password, 1 for private key
}

// CommandStep represents a single command to be executed
type CommandStep struct {
	Command string `json:"command"`
}

// RunJob 立即执行任务
func RunJob(job *model.CronJob) {
	// 检查随机延迟
	delaySeconds := 0
	if job.RandomDelay > 0 {
		// 生成 0 到 RandomDelay*60 之间的随机秒数
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		delaySeconds = r.Intn(job.RandomDelay * 60)
	}

	go func(scs int) {
		if scs > 0 {
			log.Printf("任务 [%s] 将随机延迟 %d 秒后执行...", job.Name, scs)
			time.Sleep(time.Duration(scs) * time.Second)
		}

		log.Printf("开始执行任务: %s", job.Name)

		// 解析主机信息
		var hosts []SSHInfo
		err := json.Unmarshal([]byte(job.HostInfo), &hosts)
		if err != nil {
			updateJobResult(job.ID, "失败", fmt.Sprintf("解析主机信息失败: %v", err))
			return
		}

		// 解析命令列表
		var commands []CommandStep
		err = json.Unmarshal([]byte(job.Commands), &commands)
		if err != nil {
			updateJobResult(job.ID, "失败", fmt.Sprintf("解析命令失败: %v", err))
			return
		}

		var resultLog strings.Builder
		successCount := 0
		failCount := 0

		// 遍历所有主机执行
		for _, host := range hosts {
			resultLog.WriteString(fmt.Sprintf("\n--- Host: %s (%s) ---\n", host.Hostname, host.Username))

			client := &core.SSHClient{
				Hostname:   host.Hostname,
				Port:       host.Port,
				Username:   host.Username,
				Password:   host.Password,
				PrivateKey: host.PrivateKey,
				Passphrase: host.Passphrase,
				LoginType:  host.LoginType,
			}

			// 移除 err != nil 检查，因为直接初始化不会返回错误
			// 但 RunBatchTasks 内部会调用 GenerateClient 建立连接


			// Convert CommandStep slice to string slice for RunBatchTasks
			var cmdStrings []string
			for _, cmd := range commands {
				cmdStrings = append(cmdStrings, cmd.Command)
			}

			output, err := client.RunBatchTasks(cmdStrings) // Pass string slice
			client.Close()                                  // 及时关闭

			if err != nil {
				failCount++
				resultLog.WriteString(fmt.Sprintf("执行出错: %v\nOutput:\n%s\n", err, output))
			} else {
				successCount++
				resultLog.WriteString(fmt.Sprintf("执行成功:\n%s\n", output))
			}
		}

		finalStatus := "成功"
		if failCount > 0 {
			if successCount == 0 {
				finalStatus = "失败"
			} else {
				finalStatus = "部分成功"
			}
		}

		finalLog := resultLog.String()
		updateJobResult(job.ID, finalStatus, finalLog)

		// 发送通知
		title := fmt.Sprintf("定时任务 [%s] 执行%s", job.Name, finalStatus)
		// 如果是部分成功或失败，或者策略是总是通知(这里暂定总是通知，或者可以加个NotificationStrategy字段)
		// 用户之前的需求是“及通知”，假设是全部通知，或者后续细化
		// 暂且逻辑：只要配置了通知就发
		SendNotification(title, finalLog) // Assuming SendNotification is defined elsewhere
	}(delaySeconds)
}

// updateJobResult updates the job's last run time, result, and error log in the database.
func updateJobResult(jobID uint, result, logStr string) {
	var currentJob model.CronJob
	model.DB.First(&currentJob, jobID)
	now := time.Now()
	currentJob.LastRunTime = &now
	currentJob.LastResult = result
	currentJob.ErrorLog = logStr
	model.DB.Save(&currentJob)
}

// handleJobResult is the old function, replaced by updateJobResult and the new RunJob logic.
// Keeping it for context, but it's no longer called by the new RunJob.
func handleJobResult(job *model.CronJob, runTime time.Time, result, logStr string) {
	// 重新从DB获取以避免并发覆盖
	var currentJob model.CronJob
	model.DB.First(&currentJob, job.ID)
	currentJob.LastRunTime = &runTime
	currentJob.LastResult = result
	currentJob.ErrorLog = logStr
	model.DB.Save(&currentJob)
}
