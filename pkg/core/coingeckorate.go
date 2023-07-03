package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type coingeckoResponse map[string]map[string]float64

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client for the coingecko api
// Allows to query current exchange rate of `Coin` to `Currency`
type CoingeckoRate struct {
	coin, currency string
	client         HTTPClient
	description    string
}

func getDescription(coin, currency string) string {
	cointTitle := strings.ToUpper(coin[:1]) + strings.ToLower(coin[1:])
	return fmt.Sprintf("%s to %s Exchange Rate", cointTitle, strings.ToUpper(currency))
}

func NewCoingeckoRate(coin, currency string) *CoingeckoRate {
	description := getDescription(coin, currency)
	return &CoingeckoRate{
		coin:        coin,
		currency:    currency,
		description: description,
		client:      &http.Client{},
	}
}

func NewCoingeckoRateWithHTTPClient(coin, currency string, client HTTPClient) *CoingeckoRate {
	return &CoingeckoRate{
		coin:        coin,
		currency:    currency,
		description: getDescription(coin, currency),
		client:      client,
	}
}

func (requester CoingeckoRate) Description() string {
	return requester.description
}

func (requester CoingeckoRate) Value(ctx context.Context) (float64, error) {
	// https://www.coingecko.com/en/api/documentation
	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s",
		requester.coin,
		requester.currency,
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
	rate := rateJSON[requester.coin][requester.currency]
	if err != nil {
		log.Println("CoingeckoRate.Value json error", err)
	}

	log.Println("get rate:", rate)
	return rate, err
}
