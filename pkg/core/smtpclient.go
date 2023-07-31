package core

import (
	"fmt"
	"log"
	"net/smtp"
)

type SMTPClient struct {
	from               string
	auth               smtp.Auth
	smtpHost, smtpPort string
}

func NewSMTPClient(from string, auth smtp.Auth, smtpHost, smtpPort string) SMTPClient {
	return SMTPClient{
		from:     from,
		auth:     auth,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (s SMTPClient) Send(receiver, message string) error {
	to := []string{
		receiver,
	}

	messageBytes := []byte(message)
	smtpURL := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	err := smtp.SendMail(smtpURL, s.auth, s.from, to, messageBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}

func NewAuthentication(username, password, host string) smtp.Auth {
	var auth smtp.Auth
	if password == "" {
		auth = nil
	} else {
		auth = smtp.PlainAuth("", username, password, host)
	}
	return auth
}
