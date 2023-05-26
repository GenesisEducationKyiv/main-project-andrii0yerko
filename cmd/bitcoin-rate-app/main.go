package main

import (
	"bitcoinrateapp/pkg/core"
	"errors"
	"log"
	"net/http"
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
	server := core.ExchangeRateServer{Controller: controller}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.GetRoot)
	mux.HandleFunc("/rate", server.GetRate)
	mux.HandleFunc("/subscribe", server.PostSubscribe)
	mux.HandleFunc("/sendEmails", server.PostSendEmails)

	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
