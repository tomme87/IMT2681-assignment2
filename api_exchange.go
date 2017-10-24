package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strings"
)

// handleRoot for /exchange
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		wh := Webhook{}
		err := json.NewDecoder(r.Body).Decode(&wh)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := db.Add(wh)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Fprint(w, id)
	} else {
		http.Error(w, "Not implemted", http.StatusNotImplemented)
	}
}

// handleId for /exchange/{id}
func handleId(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		wh, ok := db.Get(parts[2])
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

		ok := db.Remove(parts[2])
		if ok == false {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "Not implemted", http.StatusNotImplemented)
	}
}

// handleLatest for /exchange/latest
func handleLatest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei_latest")
}

// handleAverage for /exchange/average
func handleAverage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei_average")
}

// handleEvaluationTrigger for /exchange/evaluationtrigger
func handleEvaluationTrigger(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei_eval")
}