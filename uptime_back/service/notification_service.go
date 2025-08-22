package service

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type GmailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewGmailService(username, appPassword string) *GmailService {
	if username == "" || appPassword == "" {
		log.Println("WARNING: Gmail user/pass not set. Email notifications will be disabled.")
		return &GmailService{}
	}

	dialer := gomail.NewDialer("smtp.gmail.com", 587, username, appPassword)

	return &GmailService{
		dialer: dialer,
		from:   username,
	}
}

func (s *GmailService) SendNotification(toEmail, subject, message string) error {
	if s.dialer == nil {
		return fmt.Errorf("email service is not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"Vigil Monitor" <%s>`, s.from))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", message)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("❌ Error sending notification to %s: %v", toEmail, err)
		return err
	}

	log.Printf("✅ Notification sent to %s", toEmail)
	return nil
}
