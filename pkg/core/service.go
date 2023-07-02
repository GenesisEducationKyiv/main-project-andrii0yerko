package core

import (
	"bitcoinrateapp/pkg/rateclient"
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
	Value(ctx context.Context, coin, currency string) (rateclient.Rate, error)
}

type Sender interface {
	SendRate(receiver string, rate rateclient.Rate) error
}

type Service struct {
	receivers      Storage[string]
	rateRequester  RateRequester
	sender         Sender
	coin, currency string
}

func NewService(receivers Storage[string], rateRequester RateRequester, sender Sender) *Service {
	service := &Service{
		receivers:     receivers,
		rateRequester: rateRequester,
		sender:        sender,
		coin:          "bitcoin",
		currency:      "uah",
	}
	return service
}

func NewServiceWithDefaults(coingeckoURL, binanceURL, smtpPort, smtpHost, from, password, filename string) (*Service, error) {
	db, err := NewFileDB(filename)
	if err != nil {
		return nil, err
	}

	requester := rateclient.NewLoggingRequester(rateclient.NewCoingeckoRate(coingeckoURL))
	requesterChain := rateclient.NewRequesterChain(requester)
	requester2 := rateclient.NewLoggingRequester(rateclient.NewBinanceRate(binanceURL))
	requesterChain2 := rateclient.NewRequesterChain(requester2)
	requesterChain.SetNext(requesterChain2)

	auth := NewAuthentication(from, password, smtpHost)
	client := NewSMTPClient(from, auth, smtpHost, smtpPort)
	formatter := NewPlainEmailFormatter(from)
	sender := NewEmailSender(client, formatter)

	service := NewService(db, requesterChain, sender)
	return service, nil
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
