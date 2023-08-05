package app

import (
	"bitcoinrateapp/pkg/model"
	"bitcoinrateapp/pkg/service"
	"errors"
	"log"
	"net/http"
)

// Subscribes an email to the rate notification
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateHandler) PostSubscribe(w http.ResponseWriter, r *http.Request) {
	log.Printf("got PostSubscribe request\n")
	email := r.PostFormValue("email")

	subscriber, err := model.NewSubscriber(email)
	if errors.Is(err, model.ErrEmptyEmail) {
		return
	}

	err = e.senderService.Subscribe(subscriber)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, service.ErrIsDuplicate):
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}
