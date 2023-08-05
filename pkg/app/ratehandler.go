package app

import (
	"encoding/json"
	"log"
	"net/http"
)

// Return current exchange rate
// Accepts GET, returns rate value and StatusOK
func (e ExchangeRateHandler) GetRate(w http.ResponseWriter, _ *http.Request) {
	log.Printf("got GetRate request\n")
	value, err := e.rateService.ExchangeRate()
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

// Send email with current rate to all the subscribers
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateHandler) PostSendEmails(w http.ResponseWriter, _ *http.Request) {
	log.Printf("got PostSendEmails request\n")
	err := e.rateService.Notify()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
