package api

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strings"
)

const (
	BasePath = "/exchange"
	IdPath = "/"
	LatestPath = "/latest"
	AveragePath = "/average"
	EvaluationTriggerPath = "/evaluationtrigger"
)

var Db WebhooksStorage

// handleRoot for /exchange
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		wh := Webhook{}
		err := json.NewDecoder(r.Body).Decode(&wh)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := Db.Add(wh)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Fprint(w, id)
	} else {
		http.Error(w, "Not implemted", http.StatusNotImplemented)
	}
}

// handleId for /exchange/{id}
func HandleId(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		wh, ok := Db.Get(parts[2])
		if ok == false {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(wh)
	} else if r.Method == "DELETE" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ok := Db.Remove(parts[2])
		if ok == false {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "Not implemted", http.StatusNotImplemented)
	}
}

// handleLatest for /exchange/latest
func HandleLatest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei_latest")
}

// handleAverage for /exchange/average
func HandleAverage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei_average")
}

// handleEvaluationTrigger for /exchange/evaluationtrigger
func HandleEvaluationTrigger(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei_eval")
}