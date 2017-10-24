package main

import (
	"net/http"
	"os"
	"github.com/jasonlvhit/gocron"
	"fmt"
)

const (
	basePath = "/exchange"
	idPath = "/"
	latestPath = "/latest"
	averagePath = "/average"
	evaluationTriggerPath = "/evaluationtrigger"
)

var db WebhooksStorage

func task() {
	fmt.Println("I am runnning task.")
}

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

	gocron.Every(1).Minute().Do(task)
	go gocron.Start()

	http.ListenAndServe(":"+port, nil)
}
