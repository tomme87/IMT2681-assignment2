package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const theDiscordWebhook = "https://ptb.discordapp.com/api/webhooks/374577931869093889/_FRDtby1uvaSoy51CBo8-G_nbRPZeXaV21nTTEb1T8o3Hu23o5Ui92ExCNQ-Q8o-wgv-"

/*
WebhookInfo represents the Discord webhook data entry.
*/
type WebhookInfo struct {
	Content string `json:"content"`
}

func sendDiscordLogEntry(what string) {
	info := WebhookInfo{}
	info.Content = what + "\n"
	raw, _ := json.Marshal(info)
	resp, err := http.Post(theDiscordWebhook, "application/json", bytes.NewBuffer(raw))
	if err != nil {
		fmt.Println(err)
		fmt.Println(ioutil.ReadAll(resp.Body))
	}
}

func main() {
	for {
		text := "Heroku timer test at: " + time.Now().String()
		delay := time.Minute * 15

		sendDiscordLogEntry(text)
		time.Sleep(delay)
	}
}

/*package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/tomme87/IMT2681-assignment2/api"
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

	api.Db = &api.MongoDB{
		DatabaseURL: uri,
		DatabaseName: "exchange",
		WebhooksCollectionName: "webhooks",
		ExchangeCollectionName: "currencyrates",
	}
	api.Db.Init()

	gocron.Every(1).Minute().Do(task)
	<- gocron.Start()
}
*/