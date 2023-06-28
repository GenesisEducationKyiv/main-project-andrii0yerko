package core

import (
	"fmt"
	"log"
	"net/smtp"
)

// Simple Sender implementation around smtp package
type EmailSender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailSender(from, password, smtpHost, smtpPort string) EmailSender {
	return EmailSender{
		from:     from,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (sender EmailSender) Send(receiver string, subject, message string) error {
	// Receiver email address.
	to := []string{
		receiver,
	}

	// Message
	rfc822 := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", sender.from, receiver, subject, message)
	messageBytes := []byte(rfc822)

	// Sending email
	err := smtp.SendMail(sender.smtpHost+":"+sender.smtpPort, sender.authentication(), sender.from, to, messageBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}

func (sender EmailSender) authentication() smtp.Auth {
	// Authentication
	var auth smtp.Auth
	if sender.password == "" {
		auth = nil
	} else {
		auth = smtp.PlainAuth("", sender.from, sender.password, sender.smtpHost)
	}
	return auth
}
