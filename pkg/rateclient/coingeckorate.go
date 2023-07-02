package rateclient

import (
	"bitcoinrateapp/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type coingeckoResponse map[string]map[string]float64

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CoingeckoRate struct {
	client HTTPClient
}

func NewCoingeckoRate() *CoingeckoRate {
	return &CoingeckoRate{
		client: &http.Client{},
	}
}

func NewCoingeckoRateWithHTTPClient(client HTTPClient) *CoingeckoRate {
	return &CoingeckoRate{
		client: client,
	}
}

func (requester CoingeckoRate) Value(ctx context.Context, coin, currency string) (Rate, error) {
	// https://www.coingecko.com/en/api/documentation
	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s",
		coin,
		currency,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("CoingeckoRate.Value request error", err)
	}
	req.Header.Set("accept", "application/json")

	resp, err := requester.client.Do(req)
	if err != nil {
		log.Println("CoingeckoRate.Value api error", err)
	}

	defer resp.Body.Close()
	var rateJSON coingeckoResponse
	err = json.NewDecoder(resp.Body).Decode(&rateJSON)
	value := rateJSON[coin][currency]
	if err != nil {
		log.Println("CoingeckoRate.Value json error", err)
	}

	log.Println("get rate:", value)

	rate := model.NewExchangeRate(value, coin, currency)

	return rate, err
}
