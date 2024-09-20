package email

import (
	"auth-service/pkg/config"
	"fmt"
	"net/smtp"
)

type Emailer struct {
	SMTPServer string
	Port       string
	User       string
	Password   string
}

func NewEmailer(cfg *config.Config) *Emailer {
	return &Emailer{
		SMTPServer: "smtp.example.com",
		Port:       "587",
		User:       cfg.EmailUser,
		Password:   cfg.EmailPassword,
	}
}

func (e *Emailer) SendWarningEmail(userID, ipAddress string) error {
	from := "asilbekxolmatov2002@gmail.com"
	to := e.User
	subject := "IP Address Change Warning"
	body := fmt.Sprintf("Hello,\n\nYour IP address has changed to %s. If this wasn't you, please contact support.\n\nBest regards,\nYour Service Team", ipAddress)

	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	auth := smtp.PlainAuth("", e.User, e.Password, e.SMTPServer)

	return smtp.SendMail(fmt.Sprintf("%s:%s", e.SMTPServer, e.Port), auth, from, []string{to}, message)
}
