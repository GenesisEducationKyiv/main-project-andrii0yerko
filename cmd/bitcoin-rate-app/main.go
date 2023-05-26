package main

import (
	"bitcoinrateapp/pkg/core"
	"fmt"
)

func main() {
	var db core.Storage[string] = &core.FileDB{Filepath: "emails.dat"}
	var requester core.ValueRequester[float64] = &core.CoingeckoRate{Coin: "bitcoin", Currency: "uah"}
	var sender core.Sender = &core.EmailSender{From: "btcapp_testing@coolmail.com", Password: "", SmtpHost: "0.0.0.0", SmtpPort: "2525"}

	controller := core.Controller{
		Receivers:     db,
		RateRequester: requester,
		Sender:        sender,
	}

	controller.Subscribe("test@test.com")
	controller.Subscribe("test2@test.com")
	value, _ := controller.GetExchangeRate()
	fmt.Println(value)

	controller.Notify()
}
