package core

import (
	"fmt"
	"log"
	"net/smtp"
)

// Simple Sender implementation around smtp package
type EmailSender struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
}

func (sender EmailSender) authentication() smtp.Auth {
	// Authentication
	var auth smtp.Auth
	if sender.Password == "" {
		auth = nil
	} else {
		auth = smtp.PlainAuth("", sender.From, sender.Password, sender.SMTPHost)
	}
	return auth
}

func (sender EmailSender) Send(receiver string, subject, message string) error {
	// Receiver email address.
	to := []string{
		receiver,
	}

	// Message
	rfc822 := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", sender.From, receiver, subject, message)
	messageBytes := []byte(rfc822)

	// Sending email
	err := smtp.SendMail(sender.SMTPHost+":"+sender.SMTPPort, sender.authentication(), sender.From, to, messageBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}
