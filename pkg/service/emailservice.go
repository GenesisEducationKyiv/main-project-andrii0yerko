package service

import (
	"bitcoinrateapp/pkg/model"

	"errors"
	"log"
)

var ErrIsDuplicate = errors.New("is duplicate")

type Storage[T any] interface {
	Append(T) error
	Records() ([]T, error)
	Contains(T) bool
}

type Subscriber interface {
	Email() string
}

type Sender interface {
	SendRate(receiver string, rate model.Rate) error
}

type SenderService struct {
	receivers Storage[string]
	sender    Sender
}

func NewSenderService(receivers Storage[string], sender Sender) *SenderService {
	service := &SenderService{
		receivers: receivers,
		sender:    sender,
	}
	return service
}

func (s SenderService) Subscribe(subscriber Subscriber) error {
	receiver := subscriber.Email()
	if s.receivers.Contains(receiver) {
		return ErrIsDuplicate
	}
	return s.receivers.Append(receiver)
}

func (s SenderService) Notify(rate model.Rate) error {
	receivers, err := s.receivers.Records()
	if err != nil {
		log.Println(err)
		return err
	}

	sendErrs := make([]error, 0, len(receivers))
	for _, receiver := range receivers {
		sendErr := s.sender.SendRate(receiver, rate)
		if sendErr != nil {
			log.Println(sendErr)
			sendErrs = append(sendErrs, sendErr)
		}
	}
	if len(sendErrs) > 0 {
		err = errors.Join(sendErrs...)
		return err
	}
	return nil
}
