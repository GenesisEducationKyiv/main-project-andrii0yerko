package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

type Storage[T any] interface {
	Append(T) error
	GetRecords() ([]T, error)
}

type FileDB struct {
	filepath string
}

func (db *FileDB) GetRecords() ([]string, error) {
	var records []string
	if _, err := os.Stat(db.filepath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		log.Println("file does not exists")
		return records, nil
	}

	file, err := os.Open(db.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (db *FileDB) checkExists(value string) bool {
	records, err := db.GetRecords()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	for _, record := range records {
		if record == value {
			return true
		}
	}

	return false
}

func (db *FileDB) Append(value string) error {
	// TODO: lowercase, trim whitespaces
	if db.checkExists(value) {
		return fmt.Errorf("already exists: %s", value)
	}

	file, err := os.OpenFile(db.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	datawriter.WriteString(value + "\n")
	datawriter.Flush()
	return nil
}

func main() {
	// requester := CoingeckoRate{coin: "bitcoin", vs_currency: "uah"}
	// requester.GetValue()
	var db Storage[string] = &FileDB{"emails.dat"}
	err := db.Append("test@test.com")
	fmt.Println(err)
	err = db.Append("test2@test.com")
	fmt.Println(err)
	err = db.Append("test2@test.com")
	fmt.Println(err)
	fmt.Println(db.GetRecords())

}
