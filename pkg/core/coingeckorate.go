package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Client for the coingecko api
// Allows to query current exchange rate of `Coin` to `Currency`
type CoingeckoRate struct {
	Coin, Currency string
}

func (requester CoingeckoRate) GetDescription() string {
	cointTitle := strings.ToUpper(requester.Coin[:1]) + strings.ToLower(requester.Coin[1:])
	return fmt.Sprintf("%s to %s Exchange Rate", cointTitle, strings.ToUpper(requester.Currency))
}

func (requester CoingeckoRate) GetValue() (float64, error) {
	client := &http.Client{}
	// https://www.coingecko.com/en/api/documentation
	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s",
		requester.Coin,
		requester.Currency,
	)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		log.Println("CoingeckoRate.GetValue request error", err)
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("CoingeckoRate.GetValue api error", err)
	}

	defer resp.Body.Close()
	var rateJSON = make(map[string]map[string]float64)
	err = json.NewDecoder(resp.Body).Decode(&rateJSON)
	rate := rateJSON[requester.Coin][requester.Currency]
	if err != nil {
		log.Println("CoingeckoRate.GetValue json error", err)
	}

	log.Println("get rate:", rate)
	return rate, err
}
