package rateclient

import (
	"bitcoinrateapp/pkg/model"
	"context"
	"log"
)

type LoggingRequester struct {
	adaptee requester
}

func NewLoggingRequester(adaptee requester) *LoggingRequester {
	return &LoggingRequester{
		adaptee: adaptee,
	}
}

func (l LoggingRequester) Value(ctx context.Context, coin, currency string) (model.Rate, error) {
	rate, err := l.adaptee.Value(ctx, coin, currency)
	if err != nil {
		return nil, err
	}
	log.Printf("%T - %s/%s: %f", l.adaptee, coin, currency, rate.Value())
	return rate, nil
}
