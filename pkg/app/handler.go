package app

import (
	"bitcoinrateapp/pkg/service"
	"net/http"
)

type ExchangeRateHandler struct {
	senderService *service.SenderService
	rateService   *service.RateService
}

func NewExchangeRateHandler(
	senderService *service.SenderService,
	rateService *service.RateService,
) *ExchangeRateHandler {
	return &ExchangeRateHandler{senderService: senderService, rateService: rateService}
}

func (e ExchangeRateHandler) GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}
