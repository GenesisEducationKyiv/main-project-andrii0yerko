package service_test

import (
	"bitcoinrateapp/pkg/service"
	"bitcoinrateapp/pkg/testenv"
	"testing"
)

func TestRateServiceExchangeRate(t *testing.T) {
	rate := 100.0
	btcservice := service.NewRateService(nil, &testenv.MockRate{ExpectedRate: rate}, "bitcoin", "uah")
	actualRate, err := btcservice.ExchangeRate()
	if err != nil {
		t.Fatal(err)
	}
	if actualRate != rate {
		t.Errorf("unexpected rate: %f", actualRate)
	}
}
