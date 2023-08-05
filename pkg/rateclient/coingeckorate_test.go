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
	expectedJSON string
}

func (m *MockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	httpResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(m.expectedJSON)),
	}
	return httpResp, nil
}

func TestCoingeckoValueRequest(t *testing.T) {
	expectedRate := 1000.0

	expectedJSON := fmt.Sprintf(`{"bitcoin":{"uah":%f}}`, expectedRate)
	client := &MockHTTPClient{expectedJSON: expectedJSON}
	coingecko := rateclient.NewCoingeckoRate("https://api.coingecko.com/api/v3", client)
	actualRate, err := coingecko.Value(context.TODO(), "bitcoin", "uah")
	if err != nil {
		t.Error(err)
	}

	if math.Abs(actualRate.Value()-expectedRate) > 0.0001 {
		t.Errorf("expected %f, got %f", expectedRate, actualRate)
	}
}
