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
	SmtpHost string
	SmtpPort string
}

func (sender EmailSender) Send(receiver string, subject, message string) error {
	// Receiver email address.
	to := []string{
		receiver,
	}

	// Message
	rfc822 := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", sender.From, receiver, subject, message)
	message_bytes := []byte(rfc822)

	// Authentication
	var auth smtp.Auth
	if sender.Password == "" {
		auth = nil
	} else {
		auth = smtp.PlainAuth("", sender.From, sender.Password, sender.SmtpHost)
	}

	// Sending email
	err := smtp.SendMail(sender.SmtpHost+":"+sender.SmtpPort, auth, sender.From, to, message_bytes)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}
