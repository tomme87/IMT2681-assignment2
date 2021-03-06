package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
 BasePath is the root endpoint
 IDPath The path to get webhook by id
 LatestPath path to get latest exchange rates
 AveragePath path to get average exchange rates
 EvaluationTriggerPath to trigger all webhooks
*/
const (
	BasePath              = "/exchange"
	IDPath                = "/"
	LatestPath            = "/latest"
	AveragePath           = "/average"
	EvaluationTriggerPath = "/evaluationtrigger"
)

// Db the main Storage object
var Db Storage

// HandleRoot for /exchange  Adds a new webhook
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
			return
		}
		fmt.Fprint(w, id)
	} else {
		http.Error(w, "Not implemted", http.StatusBadRequest)
	}
}

// HandleID for /exchange/{id} Get info about webhook
func HandleID(w http.ResponseWriter, r *http.Request) {
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

// HandleLatest for /exchange/latest
func HandleLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fixer, err := Db.GetLatest(1)
		if err != nil {
			http.Error(w, "unable to get latest: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(fixer) < 1 {
			http.Error(w, "did not find latest", http.StatusInternalServerError)
			return
		}

		wh := Webhook{}
		err = json.NewDecoder(r.Body).Decode(&wh)
		if err != nil {
			http.Error(w, "unable to decode: "+err.Error(), http.StatusBadRequest)
			return
		}

		rate, err := fixer[0].GetRate(wh.BaseCurrency, wh.TargetCurrency)
		if err != nil {
			http.Error(w, "unable to convert: "+err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, rate)
	} else {
		http.Error(w, "Not implemted", http.StatusNotImplemented)
	}
}

// HandleAverage for /exchange/average
func HandleAverage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fixers, err := Db.GetLatest(3)
		if err != nil {
			http.Error(w, "unable to get latest: "+err.Error(), http.StatusInternalServerError)
			return
		}

		wh := Webhook{}
		err = json.NewDecoder(r.Body).Decode(&wh)
		if err != nil {
			http.Error(w, "unable to decode: "+err.Error(), http.StatusBadRequest)
			return
		}

		total := float32(0)
		for _, fixer := range fixers {
			rate, err := fixer.GetRate(wh.BaseCurrency, wh.TargetCurrency)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			total += rate
		}

		fmt.Fprint(w, total/float32(len(fixers)))
	} else {
		http.Error(w, "Not implemted", http.StatusNotImplemented)
	}
}

// HandleEvaluationTrigger for /exchange/evaluationtrigger This invokes all webhooks.
func HandleEvaluationTrigger(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		errs := []string{}
		for _, wh := range Db.GetAll() {
			err := wh.Invoke()
			errs = append(errs, err.Error())
		}
		json.NewEncoder(w).Encode(errs)
	} else {
		http.Error(w, "Not implemted", http.StatusNotImplemented)
	}
}
