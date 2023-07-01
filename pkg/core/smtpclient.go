package core

import (
	"fmt"
	"log"
	"net/smtp"
)

type SMTPClient struct {
	from               string
	password           string
	smtpHost, smtpPort string
}

func NewSMTPClient(from, password, smtpHost, smtpPort string) SMTPClient {
	return SMTPClient{
		from:     from,
		password: password,
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

	err := smtp.SendMail(smtpURL, s.authentication(), s.from, to, messageBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}

func (s SMTPClient) authentication() smtp.Auth {
	var auth smtp.Auth
	if s.password == "" {
		auth = nil
	} else {
		auth = smtp.PlainAuth("", s.from, s.password, s.smtpHost)
	}
	return auth
}
