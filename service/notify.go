package service

import (
	"fmt"
	"net/http"
	"net/url"
	"webssh/model"

	"gopkg.in/gomail.v2"
)

// SendNotification 发送通知
func SendNotification(subject, body string) {
	var conf model.NotificationConfig
	model.DB.First(&conf)

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
		fmt.Printf("Email send failed: %v\n", err)
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
		fmt.Printf("TG send failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func sendBark(conf model.NotificationConfig, subject, body string) {
	// Bark URL format: https://api.day.app/yourkey/title/body
	// Ensure BarkUrl ends with /
	if conf.BarkUrl == "" {
		return
	}
	
	// If user only provided the key, prepend the official API URL
	targetUrl := conf.BarkUrl
	if len(targetUrl) < 10 { // Just a rough check, key is usually longer but URL definitely is
		targetUrl = "https://api.day.app/" + targetUrl
	}
	
	// Add title and body
	// We need to encode paths
	finalUrl := fmt.Sprintf("%s/%s/%s", 
		targetUrl, 
		url.PathEscape(subject), 
		url.PathEscape(body))
		
	resp, err := http.Get(finalUrl)
	if err != nil {
		fmt.Printf("Bark send failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
}
