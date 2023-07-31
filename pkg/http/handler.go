package http

import (
	"bitcoinrateapp/pkg/core"
	"bitcoinrateapp/pkg/model"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Methods for HTTP API endpoints
type ExchangeRateHandler struct {
	Service *core.Service
}

func NewExchangeRateHandler(service *core.Service) *ExchangeRateHandler {
	return &ExchangeRateHandler{Service: service}
}

func (e ExchangeRateHandler) GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}

// Return current exchange rate
// Accepts GET, returns rate value and StatusOK
func (e ExchangeRateHandler) GetRate(w http.ResponseWriter, _ *http.Request) {
	log.Printf("got GetRate request\n")
	value, err := e.Service.ExchangeRate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResp, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Printf("Error happened in writing response. Err: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

// Subscribes an email to the rate notification
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateHandler) PostSubscribe(w http.ResponseWriter, r *http.Request) {
	log.Printf("got PostSubscribe request\n")
	email := r.PostFormValue("email")

	subscriber, err := model.NewSubscriber(email)
	if errors.Is(err, model.ErrEmptyEmail) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = e.Service.Subscribe(subscriber)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, core.ErrIsDuplicate):
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}

// Send email with current rate to all the subscribers
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateHandler) PostSendEmails(w http.ResponseWriter, _ *http.Request) {
	log.Printf("got PostSendEmails request\n")
	err := e.Service.Notify()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
