package service

import (
	"encoding/base64"
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
		var finalStatus = "尚未执行"
		var finalLog string
		var startTime = time.Now()

		defer func() {
			if r := recover(); r != nil {
				finalStatus = "失败"
				finalLog += fmt.Sprintf("\nPanic recovered: %v", r)
			}
			
			log.Printf("[Job %s] Finished with status: %s", job.Name, finalStatus)
			
			// 更新数据库
			var currentJob model.CronJob
			model.DB.First(&currentJob, job.ID)
			currentJob.LastRunTime = &startTime
			currentJob.LastResult = finalStatus
			currentJob.ErrorLog = finalLog
			model.DB.Save(&currentJob)

			// 发送通知
			title := fmt.Sprintf("定时任务 [%s] 执行%s", job.Name, finalStatus)
			SendNotification(title, finalLog)
		}()

		if scs > 0 {
			log.Printf("任务 [%s] 将随机延迟 %d 秒后执行...", job.Name, scs)
			time.Sleep(time.Duration(scs) * time.Second)
		}

		log.Printf("[Job %s] Start executing...", job.Name)

		// 1. 解析主机信息 (HostInfo)
		// Frontend 传递的是 Base64 编码的 JSON
		// 格式可能是 单个对象(Frontend) 或 数组(Legacy)
		
		var hosts []SSHInfo
		var decodedHostInfo []byte
		var err error

		// 尝试 Base64 解码
		decodedHostInfo, err = base64.StdEncoding.DecodeString(job.HostInfo)
		if err != nil {
			// 如果解码失败，尝试直接作为 JSON (兼容旧数据)
			decodedHostInfo = []byte(job.HostInfo)
		}

		// 尝试解析为单个对象
		var singleHost SSHInfo
		if err := json.Unmarshal(decodedHostInfo, &singleHost); err == nil && singleHost.Hostname != "" {
			hosts = append(hosts, singleHost)
		} else {
			// 尝试解析为数组
			if err := json.Unmarshal(decodedHostInfo, &hosts); err != nil {
				finalStatus = "失败"
				finalLog = fmt.Sprintf("解析主机信息失败: %v\nRaw: %s", err, string(decodedHostInfo))
				return
			}
		}
		
		if len(hosts) == 0 {
			finalStatus = "失败"
			finalLog = "主机列表为空"
			return
		}

		// 2. 解析命令列表 (Commands)
		// Frontend 传递的是 ["cmd1", "cmd2"] (字符串数组的JSON)
		var commands []string
		err = json.Unmarshal([]byte(job.Commands), &commands)
		if err != nil {
			// 尝试兼容旧格式 (对象数组)? 暂不，Frontend已固定
			finalStatus = "失败"
			finalLog = fmt.Sprintf("解析命令失败: %v", err)
			return
		}

		var resultLog strings.Builder
		successCount := 0
		failCount := 0

		// 3. 遍历所有主机执行
		for _, host := range hosts {
			resultLog.WriteString(fmt.Sprintf("\n--- Host: %s (%s) ---\n", host.Hostname, host.Username))
			log.Printf("[Job %s] Processing host: %s", job.Name, host.Hostname)

			client := &core.SSHClient{
				Hostname:   host.Hostname,
				Port:       host.Port,
				Username:   host.Username,
				Password:   host.Password,
				PrivateKey: host.PrivateKey,
				Passphrase: host.Passphrase,
				LoginType:  host.LoginType,
			}

			log.Printf("[Job %s] Executing commands on %s...", job.Name, host.Hostname)
			output, err := client.RunBatchTasks(commands) // commands is []string
			client.Close()
			log.Printf("[Job %s] Host %s finished. Error: %v", job.Name, host.Hostname, err)

			if err != nil {
				failCount++
				resultLog.WriteString(fmt.Sprintf("执行出错: %v\nOutput:\n%s\n", err, output))
			} else {
				successCount++
				resultLog.WriteString(fmt.Sprintf("执行成功:\n%s\n", output))
			}
		}

		if failCount > 0 {
			if successCount == 0 {
				finalStatus = "失败"
			} else {
				finalStatus = "部分成功"
			}
		} else {
			finalStatus = "成功"
		}

		finalLog = resultLog.String()
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
