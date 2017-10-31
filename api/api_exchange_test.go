package api


import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"os"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

/*
func cleanDB(db *MongoDB) {
	session, _ := mgo.Dial(db.DatabaseURL)
	session.DB(db.DatabaseName).C(db.WebhooksCollectionName).DropCollection()
	session.DB(db.DatabaseName).C(db.ExchangeCollectionName).DropCollection()
}
*/

func getTestDB() *MongoDB {
	uri := os.Getenv("MGO_TEST_URL")
	if uri == "" {
		uri = "mongodb://localhost"
	}

	TDb := &MongoDB{
		DatabaseURL: uri,
		DatabaseName: "exchange_test",
		WebhooksCollectionName: "webhooks",
		ExchangeCollectionName: "currencyrates",
	}
	TDb.Init()
	return TDb
}

func getTestWebhook() Webhook {
	wh := Webhook{
		ID: bson.NewObjectId(),
		WebhookURL: "http://test.url",
		BaseCurrency: "EUR",
		TargetCurrency: "NOK",
		MinTriggerValue: 0.2,
		MaxTriggerValue: 1.3,
	}
	return wh
}

func TestHandleRoot_ID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandleRoot))
	defer ts.Close()

	Db = getTestDB()
	Db.Init()

	//cleanDB(Db.)

	wh := getTestWebhook()

	jsonData, _ := json.Marshal(wh)

	resp, err := http.Post(ts.URL + BasePath, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error executing %s request. Error %s", http.MethodPost, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected %d, but got %d, Body: %s", http.StatusOK, resp.StatusCode, body)
		return
	}

	id, _ := ioutil.ReadAll(resp.Body)

	ts.Close()
	ts = httptest.NewServer(http.HandlerFunc(HandleID))
	defer ts.Close()

	url := ts.URL + BasePath + "/" + string(id)
	resp, err = http.Get(url)
	if err != nil {
		t.Errorf("Error executing %s request. Error %s", http.MethodGet, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected %d, but got %d, body: %s", http.StatusOK, resp.StatusCode, body)
	}

	wh2 := Webhook{}
	err = json.NewDecoder(resp.Body).Decode(&wh2)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Error executing decoding json. Error %s. Body: %s", err, body)
		return
	}

	if wh.WebhookURL != wh2.WebhookURL || wh.TargetCurrency != wh2.TargetCurrency ||
		wh.BaseCurrency != wh2.BaseCurrency || wh.MaxTriggerValue != wh2.MaxTriggerValue ||
			wh.MinTriggerValue != wh2.MinTriggerValue {
		t.Errorf("input webhook %v, do not match output webhook %v", wh, wh2)
	}
}

func TestHandleLatest(t *testing.T) {
	//ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){}))
	ts := httptest.NewServer(http.HandlerFunc(HandleLatest))
	defer ts.Close()

	Db = getTestDB()
	Db.Init()

	//wh := getTestWebhook()
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

	Db.AddCurrency(f)
	f.Date = "2017-10-27"
	Db.AddCurrency(f)

	latest := Webhook{
		BaseCurrency: "EUR",
		TargetCurrency: "NOK",
	}
	jsonData, err := json.Marshal(latest)
	if err != nil {
		t.Errorf("unable to marshall: %s", err.Error())
		return
	}

	resp, err := http.Post(ts.URL + BasePath + LatestPath, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Unable to get latest: %s", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected %d, but got %d, body: %s", http.StatusOK, resp.StatusCode, body)
		return
	}

	rates, _ := ioutil.ReadAll(resp.Body)
	rate, err := strconv.ParseFloat(string(rates), 32)
	if err != nil {
		t.Errorf("Did not get a float in return, got: %s", rates)
		return
	}

	if float32(rate) != float32(9.4865) {
		t.Errorf("Rate should be %f, got %f", float32(9.4865), float32(rate))
	}
}

