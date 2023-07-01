package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

var ErrIsDuplicate = errors.New("is duplicate")

type Storage[T any] interface {
	Append(T) error
	Records() ([]T, error)
	Contains(T) bool
}

type Rate interface {
	Value() float64
	Coin() string
	Currency() string
}

type RateRequester interface {
	Value(ctx context.Context, coin, currency string) (Rate, error)
}

type Sender interface {
	Send(receiver string, subject string, message string) error
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

func NewServiceWithDefaults(smtpPort, smtpHost, from, password, filename string) (*Service, error) {
	db, err := NewFileDB(filename)
	if err != nil {
		return nil, err
	}

	requester := NewCoingeckoRate("bitcoin", "uah")
	sender := NewEmailSender(from, password, smtpHost, smtpPort)

	service := NewService(db, requester, sender)
	return service, nil
}

func (s Service) ExchangeRate() (float64, error) {
	rate, err := s.rateRequester.Value(context.TODO(), s.coin, s.currency)
	if err != nil {
		return 0, err
	}
	return rate.Value(), nil
}

func (s Service) Subscribe(receiver string) error {
	receiver = strings.ToLower(strings.TrimSpace(receiver))
	if s.receivers.Contains(receiver) {
		return ErrIsDuplicate
	}
	return s.receivers.Append(receiver)
}

func (s Service) Notify() error {
	value, err := s.rateRequester.Value(context.TODO(), s.coin, s.currency)
	if err != nil {
		log.Println(err)
		return err
	}
	subject := getDescription(s.coin, s.currency)
	message := fmt.Sprintf("%f", value.Value())

	receivers, err := s.receivers.Records()
	if err != nil {
		log.Println(err)
		return err
	}

	sendErrs := make([]error, 0, len(receivers))
	for _, receiver := range receivers {
		sendErr := s.sender.Send(receiver, subject, message)
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

func getDescription(coin, currency string) string {
	cointTitle := strings.ToUpper(coin[:1]) + strings.ToLower(coin[1:])
	return fmt.Sprintf("%s to %s Exchange Rate", cointTitle, strings.ToUpper(currency))
}
