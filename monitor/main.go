package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/tomme87/IMT2681-assignment2/api"
	"os"
	"strings"
)

func updateTicker() {
	fmt.Println("Fetching new exchange rates.")
	fixer, err := api.NewFixer()
	if err != nil {
		fmt.Printf("Error fetching rates: %s\n", err.Error())
		return
	}

	err = api.Db.AddCurrency(*fixer)
	if err != nil {
		fmt.Printf("Error adding Currency to db: %s\n", err.Error())
		if !strings.HasPrefix(err.Error(), "E11000") { // Return unless it's a duplicate key error.
			return
		}
	}

	webhooks := api.Db.GetAll()
	for _, wh := range webhooks {
		cRate, err := fixer.GetRate(wh.BaseCurrency, wh.TargetCurrency)
		if err != nil {
			fmt.Printf("Error getting rate for %s: %s\n", wh.ID.Hex(), err.Error())
			continue
		}

		if cRate < wh.MinTriggerValue || cRate > wh.MaxTriggerValue {
			fmt.Printf("Invoking %s", wh.WebhookURL)
			wh.CurrentRate = cRate
			err := wh.Invoke()
			if err != nil {
				fmt.Printf("Error invoking %s (%s): %s\n", wh.WebhookURL, wh.ID.Hex(), err.Error())
				continue
			}
		}
	}
	fmt.Println("done")
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