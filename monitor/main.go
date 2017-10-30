package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/tomme87/IMT2681-assignment2/api"
	"os"
)

func updateTicker() {
	fmt.Println("Fetching new exchange rates.")
	fixer, err := api.NewFixer()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = api.Db.AddCurrency(*fixer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}



func main() {
	uri := os.Getenv("MGO_URL")
	if uri == "" {
		uri = "mongodb://localhost"
	}

	api.Db = &api.MongoDB{
		DatabaseURL: uri,
		DatabaseName: "exchange",
		WebhooksCollectionName: "webhooks",
		ExchangeCollectionName: "currencyrates",
	}
	api.Db.Init()

	updateTicker()
	gocron.Every(1).Day().At("17:00").Do(updateTicker)
	<- gocron.Start()
}