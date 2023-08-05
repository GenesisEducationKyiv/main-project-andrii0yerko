package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Return current exchange rate
// Accepts GET, returns rate value and StatusOK
func (e ExchangeRateHandler) GetRate(w http.ResponseWriter, _ *http.Request) {
	e.logger.Debug("got GetRate request")
	value, err := e.rateService.ExchangeRate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResp, err := json.Marshal(value)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error happened in JSON marshal. Err: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error happened in writing response. Err: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e.logger.Info("rate returned successfully")
	w.Header().Set("Content-Type", "application/json")
}

// Send email with current rate to all the subscribers
// Accepts POST, returns StatusOK if the email was not subscribed before and StatusConflict otherwise
func (e ExchangeRateHandler) PostSendEmails(w http.ResponseWriter, _ *http.Request) {
	e.logger.Debug("got PostSendEmails request\n")
	err := e.rateService.Notify()
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error happened in sending emails. Err: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e.logger.Info("emails sent successfully")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
