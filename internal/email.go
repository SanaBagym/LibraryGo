package internal

import (
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	Username     string
	Password     string
	FromAddress  string
}

var Mailer *EmailConfig

func InitializeMailer() {
	Mailer = &EmailConfig{
		SMTPHost:    "smtp.example.com",
		SMTPPort:    587,
		Username:    "your_email@example.com",
		Password:    "your_password",
		FromAddress: "noreply@example.com",
	}
}

func (ec *EmailConfig) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", ec.Username, ec.Password, ec.SMTPHost)

	emailBody := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	addr := fmt.Sprintf("%s:%d", ec.SMTPHost, ec.SMTPPort)
	if err := smtp.SendMail(addr, auth, ec.FromAddress, []string{to}, []byte(emailBody)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
