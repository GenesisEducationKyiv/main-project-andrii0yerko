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
	Receivers     Storage[string]
	RateRequester ValueRequester[float64]
	Sender        Sender
}

func (c Controller) GetExchangeRate() (float64, error) {
	return c.RateRequester.GetValue()
}

func (c Controller) Subscribe(receiver string) error {
	receiver = strings.ToLower(strings.TrimSpace(receiver))
	return c.Receivers.Append(receiver)
}

func (c Controller) Notify() error {
	value, err := c.GetExchangeRate()
	if err != nil {
		log.Println(err)
		return err
	}
	subject := c.RateRequester.GetDescription()
	message := fmt.Sprintf("%f", value)

	receivers, err := c.Receivers.GetRecords()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, receiver := range receivers {
		c.Sender.Send(receiver, subject, message)
	}

	return nil
}
