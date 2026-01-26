package utils

import (
	"net/smtp"
	"os"
)

func SendMail(to string, subject string, body string) error {

	host := os.Getenv("SMTP_HOST")
	err := smtp.SendMail(host, nil, "no-reply@school.com", []string{to}, []byte(body))
	if err != nil {
		return err
	}
	return nil
}
