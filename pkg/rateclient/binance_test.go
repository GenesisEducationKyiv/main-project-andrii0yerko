package rateclient_test

import (
	"bitcoinrateapp/pkg/rateclient"
	"context"
	"fmt"
	"math"
	"testing"
)

func TestBinanceValueRequest(t *testing.T) {
	expectedRate := 1000.0

	expectedJSON := fmt.Sprintf(`{"symbol":"BTCUAH","price":"%f"}`, expectedRate)
	client := &MockHTTPClient{expectedJSON: expectedJSON}
	binance := rateclient.NewBinanceRateWithHTTPClient("https://www.binance.com/api/v3", client)
	actualRate, err := binance.Value(context.TODO(), "bitcoin", "uah")
	if err != nil {
		t.Error(err)
	}

	if math.Abs(actualRate.Value()-expectedRate) > 0.0001 {
		t.Errorf("expected %f, got %f", expectedRate, actualRate)
	}
}
