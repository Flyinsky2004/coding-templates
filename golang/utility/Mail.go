/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package utility

import (
	"example/config"

	"gopkg.in/gomail.v2"
)

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Config.Mail.SenderName)
	m.SetHeader("From", config.Config.Mail.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.Config.Mail.SmtpHost, config.Config.Mail.SmtpPort, config.Config.Mail.Username, config.Config.Mail.Password)

	return d.DialAndSend(m)
}
