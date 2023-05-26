package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type CoingeckoRate struct {
	Coin, Currency string
}

func (requester CoingeckoRate) GetDescription() string {
	return fmt.Sprintf("%s to %s Exchange Rate", strings.Title(requester.Coin), strings.ToUpper(requester.Currency))
}

func (requester CoingeckoRate) GetValue() (float64, error) {
	client := &http.Client{}
	// https://www.coingecko.com/en/api/documentation
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", requester.Coin, requester.Currency)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var rateJson struct {
		Bitcoin struct{ Uah float64 }
	}
	err = json.NewDecoder(resp.Body).Decode(&rateJson)
	rate := rateJson.Bitcoin.Uah

	if err != nil {
		log.Fatal(err)
	}

	log.Println("get rate:", rate)
	return rate, err
}
