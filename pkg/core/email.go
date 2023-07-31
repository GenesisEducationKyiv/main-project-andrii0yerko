package core

import (
	"bitcoinrateapp/pkg/rateclient"
	"log"
)

type Formatter interface {
	Format(user string, rate rateclient.Rate) string
}

type Client interface {
	Send(receiver, message string) error
}

type EmailSender struct {
	client    Client
	formatter Formatter
}

func NewEmailSender(client Client, formatter Formatter) EmailSender {
	return EmailSender{
		client:    client,
		formatter: formatter,
	}
}

func (s EmailSender) SendRate(receiver string, rate rateclient.Rate) error {
	rfc822 := s.formatter.Format(receiver, rate)
	err := s.client.Send(receiver, rfc822)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}
