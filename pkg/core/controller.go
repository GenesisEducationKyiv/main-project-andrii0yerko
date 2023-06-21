package core

import (
	"fmt"
	"log"
	"strings"
)

// An abstract storage which allows to read and add values
type Storage[T any] interface {
	Append(T) error
	GetRecords() ([]T, error)
}

// Abstract requester which allows to extract a specific value, and its description
type ValueRequester[T any] interface {
	GetValue() (T, error)
	GetDescription() string
}

// Defines behavior of sending data for the users
type Sender interface {
	Send(receiver string, subject string, message string) error
}

// handles main logic of the App.
// responsible for providing access to the aggregated core objects
// and for setting up their interaction as well
type Controller struct {
	receivers     Storage[string]
	rateRequester ValueRequester[float64]
	sender        Sender
}

func NewController(smtpPort, smtpHost, from, password, filename string) Controller {
	var db Storage[string] = NewFileDB(filename)
	var requester ValueRequester[float64] = NewCoingeckoRate("bitcoin", "uah")
	var sender Sender = NewEmailSender(from, password, smtpHost, smtpPort)

	controller := Controller{
		receivers:     db,
		rateRequester: requester,
		sender:        sender,
	}
	return controller
}

func (c Controller) GetExchangeRate() (float64, error) {
	return c.rateRequester.GetValue()
}

func (c Controller) Subscribe(receiver string) error {
	receiver = strings.ToLower(strings.TrimSpace(receiver))
	return c.receivers.Append(receiver)
}

func (c Controller) Notify() error {
	value, err := c.GetExchangeRate()
	if err != nil {
		log.Println(err)
		return err
	}
	subject := c.rateRequester.GetDescription()
	message := fmt.Sprintf("%f", value)

	receivers, err := c.receivers.GetRecords()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, receiver := range receivers {
		sendErr := c.sender.Send(receiver, subject, message)
		if sendErr != nil {
			log.Println(sendErr)
			err = sendErr
		}
	}

	return err
}
