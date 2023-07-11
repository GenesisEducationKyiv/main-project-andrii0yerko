package service

import (
	"bitcoinrateapp/pkg/model"

	"context"
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

type RateRequester interface {
	Value(ctx context.Context, coin, currency string) (model.Rate, error)
}

type Sender interface {
	SendRate(receiver string, rate model.Rate) error
}

type Service struct {
	receivers      Storage[string]
	rateRequester  RateRequester
	sender         Sender
	coin, currency string
}

func NewService(receivers Storage[string], rateRequester RateRequester, sender Sender, coin, currency string) *Service {
	service := &Service{
		receivers:     receivers,
		rateRequester: rateRequester,
		sender:        sender,
		coin:          coin,
		currency:      currency,
	}
	return service
}

func (s Service) ExchangeRate() (float64, error) {
	rate, err := s.rateRequester.Value(context.TODO(), s.coin, s.currency)
	if err != nil {
		return 0, err
	}
	return rate.Value(), nil
}

func (s Service) Subscribe(subscriber Subscriber) error {
	receiver := subscriber.Email()
	if s.receivers.Contains(receiver) {
		return ErrIsDuplicate
	}
	return s.receivers.Append(receiver)
}

func (s Service) Notify() error {
	rate, err := s.rateRequester.Value(context.TODO(), s.coin, s.currency)
	if err != nil {
		log.Println(err)
		return err
	}

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
