package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ValueRequester interface {
	GetValue() (float64, error)
}

type CoingeckoRate struct {
	coin, vs_currency string
}

func (requester CoingeckoRate) GetValue() (float64, error) {
	client := &http.Client{}
	// https://www.coingecko.com/en/api/documentation
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", requester.coin, requester.vs_currency)

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

func main() {
	requester := CoingeckoRate{coin: "bitcoin", vs_currency: "uah"}
	requester.GetValue()
}
