package app

import (
	"bitcoinrateapp/pkg/service"
	"net/http"
)

type Logger interface {
	Debug(message string)
	Info(message string)
	Error(message string)
}

type ExchangeRateHandler struct {
	senderService *service.SenderService
	rateService   *service.RateService
	logger        Logger
}

func NewExchangeRateHandler(
	senderService *service.SenderService,
	rateService *service.RateService,
	logger Logger,
) *ExchangeRateHandler {
	return &ExchangeRateHandler{senderService: senderService, rateService: rateService, logger: logger}
}

func (e ExchangeRateHandler) GetRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}
