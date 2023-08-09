package service

import (
	"bitcoinrateapp/pkg/model"
	"log"

	"context"
)

type RateRequester interface {
	Value(ctx context.Context, coin, currency string) (model.Rate, error)
}

type SenderFacade interface {
	Notify(rate model.Rate) error
}

type RateService struct {
	rateRequester  RateRequester
	coin, currency string
	senderClient   SenderFacade
}

func NewRateService(sender SenderFacade, rateRequester RateRequester, coin, currency string) *RateService {
	service := &RateService{
		rateRequester: rateRequester,
		coin:          coin,
		currency:      currency,
		senderClient:  sender,
	}
	return service
}

func (s RateService) ExchangeRate() (float64, error) {
	rate, err := s.rateRequester.Value(context.TODO(), s.coin, s.currency)
	if err != nil {
		return 0, err
	}
	return rate.Value(), nil
}

func (s RateService) Notify() error {
	rate, err := s.rateRequester.Value(context.TODO(), s.coin, s.currency)
	if err != nil {
		log.Println(err)
		return err
	}
	return s.senderClient.Notify(rate)
}
