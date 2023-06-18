package core

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Methods for HTTP API endpoints
type ExchangeRateServer struct {
	Controller Controller
}

func (e ExchangeRateServer) GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}

// Return current exchange rate
// Accepts GET, returns rate value and StatusOK
func (e ExchangeRateServer) GetRate(w http.ResponseWriter, r *http.Request) {
	if len(r.Method) > 0 && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("got GetRate request\n")
	value, err := e.Controller.GetExchangeRate()
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
}

// Subscribes an email to the rate notification
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateServer) PostSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("got PostSubscribe request\n")
	email := r.PostFormValue("email")
	if email == "" {
		return
	}
	err := e.Controller.Subscribe(email)

	if err != nil {
		switch {
		case errors.Is(err, ErrIsDuplicate):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Send email with current rate to all the subscribers
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateServer) PostSendEmails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Printf("got PostSendEmails request\n")
	err := e.Controller.Notify()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
