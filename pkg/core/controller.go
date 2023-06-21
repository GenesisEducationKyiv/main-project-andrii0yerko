package core

import (
	"fmt"
	"log"
	"strings"
)

// handles main logic of the App.
// responsible for providing access to the aggregated core objects
// and for setting up their interaction as well
type Controller struct {
	receivers     Storage[string]
	rateRequester ValueRequester[float64]
	sender        Sender
}

func NewController(smtpPort, smtpHost, from, password, filename string) Controller {
	var db Storage[string] = &FileDB{Filepath: filename}
	var requester ValueRequester[float64] = &CoingeckoRate{Coin: "bitcoin", Currency: "uah"}
	var sender Sender = &EmailSender{From: from, Password: password, SMTPHost: smtpHost, SMTPPort: smtpPort}

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
