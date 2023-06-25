package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

// An abstract storage which allows to read and add values
type Storage[T any] interface {
	Append(T) error
	Records() ([]T, error)
}

// Abstract requester which allows to extract a specific value, and its description
type ValueRequester[T any] interface {
	Value(context.Context) (T, error)
	Description() string
}

// Defines behavior of sending data for the users
type Sender interface {
	Send(receiver string, subject string, message string) error
}

// handles main logic of the App.
// responsible for providing access to the aggregated core objects
// and for setting up their interaction as well
type Service struct {
	receivers     Storage[string]
	rateRequester ValueRequester[float64]
	sender        Sender
}

func NewService(receivers Storage[string], rateRequester ValueRequester[float64], sender Sender) *Service {
	service := &Service{
		receivers:     receivers,
		rateRequester: rateRequester,
		sender:        sender,
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
	return s.rateRequester.Value(context.TODO())
}

func (s Service) Subscribe(receiver string) error {
	receiver = strings.ToLower(strings.TrimSpace(receiver))
	return s.receivers.Append(receiver)
}

func (s Service) Notify() error {
	value, err := s.ExchangeRate()
	if err != nil {
		log.Println(err)
		return err
	}
	subject := s.rateRequester.Description()
	message := fmt.Sprintf("%f", value)

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
