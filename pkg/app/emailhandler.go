package app

import (
	"bitcoinrateapp/pkg/model"
	"bitcoinrateapp/pkg/service"
	"errors"
	"fmt"
	"net/http"
)

// Subscribes an email to the rate notification
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateHandler) PostSubscribe(w http.ResponseWriter, r *http.Request) {
	e.logger.Debug("got PostSubscribe request")
	email := r.PostFormValue("email")

	subscriber, err := model.NewSubscriber(email)
	if errors.Is(err, model.ErrEmptyEmail) {
		return
	}

	err = e.senderService.Subscribe(subscriber)
	switch {
	case err == nil:
		e.logger.Info("email subscribed successfully")
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, service.ErrIsDuplicate):
		e.logger.Info("email is already subscribed")
		w.WriteHeader(http.StatusConflict)
	default:
		e.logger.Error(fmt.Sprintf("Error happened in subscribing email. Err: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}
