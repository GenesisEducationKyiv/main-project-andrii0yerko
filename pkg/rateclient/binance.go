package rateclient

import (
	"bitcoinrateapp/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type binanceResponse struct {
	Price float64 `json:"price,string"`
}

type BinanceRate struct {
	client  HTTPClient
	baseURL string
}

func NewBinanceRate(binanceURL string, client HTTPClient) *BinanceRate {
	return &BinanceRate{
		client:  client,
		baseURL: binanceURL,
	}
}

func (b BinanceRate) Value(ctx context.Context, coin, currency string) (Rate, error) {
	ticker := b.coinToTicker(coin)
	if ticker == "" {
		return nil, fmt.Errorf("unknown coin %s", coin)
	}
	url := fmt.Sprintf(
		"%s/ticker/price?symbol=%s%s",
		b.baseURL,
		ticker,
		strings.ToUpper(currency),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var rateJSON binanceResponse
	err = json.NewDecoder(resp.Body).Decode(&rateJSON)
	value := rateJSON.Price
	if err != nil {
		return nil, err
	}

	rate := model.NewExchangeRate(value, coin, currency)

	return rate, err
}

func (b BinanceRate) coinToTicker(coin string) string {
	var mapping = map[string]string{"bitcoin": "BTC"}
	return mapping[coin]
}
