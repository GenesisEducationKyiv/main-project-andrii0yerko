package rateclient

import (
	"bitcoinrateapp/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type coingeckoResponse map[string]map[string]float64

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CoingeckoRate struct {
	client  HTTPClient
	baseURL string
}

func NewCoingeckoRate(coingeckoURL string, client HTTPClient) *CoingeckoRate {
	return &CoingeckoRate{
		client:  client,
		baseURL: coingeckoURL,
	}
}

func (c CoingeckoRate) Value(ctx context.Context, coin, currency string) (model.Rate, error) {
	// https://www.coingecko.com/en/api/documentation
	url := fmt.Sprintf(
		"%s/simple/price?ids=%s&vs_currencies=%s",
		c.baseURL,
		coin,
		currency,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var rateJSON coingeckoResponse
	err = json.NewDecoder(resp.Body).Decode(&rateJSON)
	value := rateJSON[coin][currency]
	if err != nil {
		return nil, err
	}

	rate := model.NewExchangeRate(value, coin, currency)

	return rate, err
}
