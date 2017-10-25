package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"IMT2681-assignment2/api"
	"os"
)

func task() {
	fmt.Println("I am runnning worker task.")
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

	api.Db = &api.WebhooksMongoDB{
		DatabaseURL: uri,
		DatabaseName: "exchange",
		WebhooksCollectionName: "webhooks",
	}
	api.Db.Init()

	gocron.Every(1).Minute().Do(task)
	<- gocron.Start()
}


