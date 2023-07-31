package rateclient_test

import (
	"bitcoinrateapp/pkg/rateclient"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
)

type MockHTTPClient struct {
	expectedRate float64
}

func (m *MockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	expectedJSON := fmt.Sprintf(`{"bitcoin":{"uah":%f}}`, m.expectedRate)
	httpResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(expectedJSON)),
	}
	return httpResp, nil
}

func TestValueRequest(t *testing.T) {
	expectedRate := 1000.0

	client := &MockHTTPClient{expectedRate: expectedRate}
	coingecko := rateclient.NewCoingeckoRateWithHTTPClient("https://api.coingecko.com/api/v3", client)
	actualRate, err := coingecko.Value(context.TODO(), "bitcoin", "uah")
	if err != nil {
		t.Error(err)
	}

	if math.Abs(actualRate.Value()-expectedRate) > 0.0001 {
		t.Errorf("expected %f, got %f", expectedRate, actualRate)
	}
}
