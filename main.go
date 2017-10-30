package main

import (
	"net/http"
	"os"
	"github.com/tomme87/IMT2681-assignment2/api"
)

func main() {
	port := os.Getenv("PORT") // Get port from environment variable. Needed to deploy on heruko.
	if port == "" {
		port = "8080" // Default to port 8080
	}

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

	http.HandleFunc(api.BasePath, api.HandleRoot)
	http.HandleFunc(api.BasePath + api.IdPath, api.HandleId)
	http.HandleFunc(api.BasePath + api.LatestPath, api.HandleLatest)
	http.HandleFunc(api.BasePath + api.AveragePath, api.HandleAverage)
	http.HandleFunc(api.BasePath + api.EvaluationTriggerPath, api.HandleEvaluationTrigger)
	http.HandleFunc(api.BasePath + "/update", api.HandleUpdateTicker)

	http.ListenAndServe(":"+port, nil)
}
