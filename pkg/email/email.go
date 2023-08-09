package email

import (
	"bitcoinrateapp/pkg/model"
	"log"
)

type Formatter interface {
	Format(user string, rate model.Rate) string
}

type Client interface {
	Send(receiver, message string) error
}

type Sender struct {
	client    Client
	formatter Formatter
}

func NewSender(client Client, formatter Formatter) Sender {
	return Sender{
		client:    client,
		formatter: formatter,
	}
}

func (s Sender) SendRate(receiver string, rate model.Rate) error {
	rfc822 := s.formatter.Format(receiver, rate)
	err := s.client.Send(receiver, rfc822)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Email Sent Successfully!")
	return nil
}
