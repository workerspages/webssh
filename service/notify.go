package service

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"webssh/model"

	"gopkg.in/gomail.v2"
)

// SendNotification 发送通知
func SendNotification(subject, body string) {
	var conf model.NotificationConfig
	if err := model.DB.First(&conf).Error; err != nil {
		log.Printf("[Notification] Failed to load config: %v", err)
		return
	}

	log.Printf("[Notification] Attempting to send '%s'. Email: %v, TG: %v, Bark: %v", subject, conf.EnableEmail, conf.EnableTg, conf.EnableBark)

	if conf.EnableEmail {
		go sendEmail(conf, subject, body)
	}
	if conf.EnableTg {
		go sendTelegram(conf, subject, body)
	}
	if conf.EnableBark {
		go sendBark(conf, subject, body)
	}
}

func sendEmail(conf model.NotificationConfig, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.EmailUser)
	m.SetHeader("To", conf.EmailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(conf.EmailHost, conf.EmailPort, conf.EmailUser, conf.EmailPass)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("[Notification] Email send failed: %v", err)
	} else {
		log.Println("[Notification] Email sent successfully")
	}
}

func sendTelegram(conf model.NotificationConfig, subject, body string) {
	msg := fmt.Sprintf("%s\n\n%s", subject, body)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", conf.TgBotToken)
	
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {conf.TgChatID},
		"text":    {msg},
	})
	if err != nil {
		log.Printf("[Notification] TG send failed: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Println("[Notification] TG sent successfully")
}

func sendBark(conf model.NotificationConfig, subject, body string) {
	// Bark URL format: https://api.day.app/yourkey/title/body
	// Ensure BarkUrl ends with /
	if conf.BarkUrl == "" {
		log.Println("[Notification] Bark skipped: URL is empty")
		return
	}
	
	// If user only provided the key, prepend the official API URL
	targetUrl := conf.BarkUrl
	if len(targetUrl) < 10 { // Just a rough check
		targetUrl = "https://api.day.app/" + targetUrl
	}
	
	// Add title and body
	finalUrl := fmt.Sprintf("%s/%s/%s", 
		strings.TrimRight(targetUrl, "/"), // Ensure clean join 
		url.PathEscape(subject), 
		url.PathEscape(body))
		
	log.Printf("[Notification] Sending Bark to: %s", finalUrl)

	resp, err := http.Get(finalUrl)
	if err != nil {
		log.Printf("[Notification] Bark send failed: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		log.Printf("[Notification] Bark returned status: %d", resp.StatusCode)
	} else {
		log.Println("[Notification] Bark sent successfully")
	}
}
