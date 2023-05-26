package main

import (
	"bitcoinrateapp/pkg/core"
	"fmt"
)

func main() {
	// requester := CoingeckoRate{coin: "bitcoin", vs_currency: "uah"}
	// requester.GetValue()
	var db core.Storage[string] = &core.FileDB{Filepath: "emails.dat"}
	err := db.Append("test@test.com")
	fmt.Println(err)
	err = db.Append("test2@test.com")
	fmt.Println(err)
	err = db.Append("test2@test.com")
	fmt.Println(err)
	fmt.Println(db.GetRecords())

}
