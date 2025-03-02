package utility

import (
	"example/config"

	"gopkg.in/gomail.v2"
)

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.MailSenderName)
	m.SetHeader("From", config.MailUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.MailSmtpHost, config.MailSmtpPort, config.MailUsername, config.MailPassword)

	return d.DialAndSend(m)
}
