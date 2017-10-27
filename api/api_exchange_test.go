package api

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

func TestHandleRoot(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(HandleRoot))
	defer ts.Close()


	Db = &MongoDB{
		DatabaseURL: "mongodb://tomme:twick493@192.168.2.60/WebhooksDB",
		DatabaseName: "WebhooksDB",
		WebhooksCollectionName: "webhooks",
		ExchangeCollectionName: "currencyrates",
	}

	Db.Init()

	wh := Webhook{
		WebhookURL: "http://test.url",
		BaseCurrency: "NOK",
		TargetCurrency: "USD",
		MinTriggerValue: 10,
		MaxTriggerValue: 20,
	}
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
	ts = httptest.NewServer(http.HandlerFunc(HandleId))
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