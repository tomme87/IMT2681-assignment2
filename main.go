package main

import (
	"net/http"
	"os"
)

const (
	basePath = "/exchange"
	idPath = "/"
	latestPath = "/latest"
	averagePath = "/average"
	evaluationTriggerPath = "/evaluationtrigger"
)

var db WebhooksStorage

func main() {
	port := os.Getenv("PORT") // Get port from environment variable. Needed to deploy on heruko.
	if port == "" {
		port = "8080" // Default to port 8080
	}

	uri := os.Getenv("MGO_URL")
	if uri == "" {
		uri = "mongodb://localhost"
	}

	db = &WebhooksMongoDB{
		DatabaseURL: uri,
		DatabaseName: "exchange",
		WebhooksCollectionName: "webhooks",
	}
	db.Init()

	http.HandleFunc(basePath, handleRoot)
	http.HandleFunc(basePath + idPath, handleId)
	http.HandleFunc(basePath + latestPath, handleLatest)
	http.HandleFunc(basePath + averagePath, handleAverage)
	http.HandleFunc(basePath + evaluationTriggerPath, handleEvaluationTrigger)

	http.ListenAndServe(":8080", nil)
}
