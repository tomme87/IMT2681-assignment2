package api

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
)

// TestFixer_GetRate to test getting the rate from fixer data.
func TestFixer_GetRate(t *testing.T) {
	f := Fixer{
		Base: "EUR",
		Date: "2017-10-26",
		Rates: map[string]float32{
			"AUD":1.5248,"BGN":1.9558,"BRL":3.803,"CAD":1.5041,"CHF":1.1678,
			"CNY":7.8003,"CZK":25.589,"DKK":7.4432,"GBP":0.8901,"HKD":9.1701,
			"HRK":7.5155,"HUF":310.32,"IDR":15982.0,"ILS":4.1343,"INR":76.23,
			"JPY":133.75,"KRW":1320.4,"MXN":22.368,"MYR":4.9762,"NOK":9.4865,
			"NZD":1.7118,"PHP":60.939,"PLN":4.235,"RON":4.5983,"RUB":67.76,
			"SEK":9.7218,"SGD":1.601,"THB":38.973,"TRY":4.4338,"USD":1.1753,"ZAR":16.739},
	}

	rate, err := f.GetRate("EUR", "NOK")
	if err != nil {
		t.Errorf("EUR and NOK should be valid. got: %s", err.Error())
		return
	}

	if rate != 9.4865 {
		t.Errorf("Expected rate 9.4865, got %f", rate)
	}

	rate, err = f.GetRate("NOK", "EUR")
	if err != nil {
		t.Errorf("EUR and NOK should be valid. got: %s", err.Error())
		return
	}

	if rate != 1 / float32(9.4865) {
		t.Errorf("Expected rate %f, got %f", 1 / float32(9.4865), rate)
	}

	rate, err = f.GetRate("SEK", "NOK")
	if err != nil {
		t.Errorf("EUR and NOK should be valid. got: %s", err.Error())
		return
	}

	if rate != float32(9.4865) * (1 / float32(9.7218)) {
		t.Errorf("Expected rate %f, got %f", float32(9.4865) * (1 / float32(9.7218)), rate)
	}
}

// handleTestGetLatest to simulate fixer API
func handleTestGetLatest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json := `{"base":"EUR","date":"2017-10-26","rates":{"AUD":1.5248,"BGN":1.9558,"BRL":3.803,"CAD":1.5041,
			"CHF":1.1678,"CNY":7.8003,"CZK":25.589,"DKK":7.4432,"GBP":0.8901,"HKD":9.1701,"HRK":7.5155,"HUF":310.32,
			"IDR":15982.0,"ILS":4.1343,"INR":76.23,"JPY":133.75,"KRW":1320.4,"MXN":22.368,"MYR":4.9762,"NOK":9.4865,
			"NZD":1.7118,"PHP":60.939,"PLN":4.235,"RON":4.5983,"RUB":67.76,"SEK":9.7218,"SGD":1.601,"THB":38.973,
			"TRY":4.4338,"USD":1.1753,"ZAR":16.739}}`
		fmt.Fprint(w, json)
	default:
		http.Error(w, "not implemented", http.StatusBadRequest)
	}
}

// TestNewFixer to test getting fixer data. using a test server
func TestNewFixer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleTestGetLatest))
	defer ts.Close()

	f, err := NewFixer(ts.URL)
	if err != nil {
		t.Errorf("Error getting fixer: %s", err.Error())
	}

	if f.Base != "EUR" {
		t.Errorf("Expected base EUR, got %s", f.Base)
	}

	if f.Date != "2017-10-26" {
		t.Errorf("Expected date 2017-10-26, got %s", f.Date)
	}

	if f.Rates["NOK"] != float32(9.4865) {
		t.Errorf("Expected date 2017-10-26, got %s", f.Rates["NOK"])
	}
}