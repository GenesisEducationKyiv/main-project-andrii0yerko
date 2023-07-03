package core_test

import (
	"bitcoinrateapp/pkg/core"
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
	coingecko := core.NewCoingeckoRateWithHTTPClient("bitcoin", "uah", client)
	actualRate, err := coingecko.Value(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if math.Abs(actualRate-expectedRate) > 0.0001 {
		t.Errorf("expected %f, got %f", expectedRate, actualRate)
	}
}
