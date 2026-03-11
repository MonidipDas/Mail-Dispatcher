package smtp

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func SendMail(config SMTPConfig, to string, subject string, body string) error {

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	msg := []byte(
		"Subject: " + subject + "\r\n" +
			"\r\n" + body + "\r\n")

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	err := smtp.SendMail(
		addr,
		auth,
		config.Username,
		[]string{to},
		msg,
	)

	return err
}
