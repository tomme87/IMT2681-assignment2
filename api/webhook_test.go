package api

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
)

// TestWebhook_Validate tests the validate function
func TestWebhook_Validate(t *testing.T) {
	wh := Webhook{
		ID: bson.NewObjectId(),
		WebhookURL: "http://test.url",
		BaseCurrency: "USD",
		TargetCurrency: "NOK",
		MinTriggerValue: 0.2,
		MaxTriggerValue: 1.3,
	}

	err := wh.Validate()
	if err != nil {
		t.Errorf("Unable to validate: %s",err.Error())
	}
}

// TestWebhook_invoke test invoking webhook.
func TestWebhook_Invoke(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST", http.StatusBadRequest)
		} else {
			wh := Webhook{}
			err := json.NewDecoder(r.Body).Decode(&wh)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if wh.CurrentRate != float32(9.5348) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
	}))
	defer ts.Close()

	wh := Webhook{
		WebhookURL: ts.URL,
		BaseCurrency: "USD",
		TargetCurrency: "NOK",
		CurrentRate: 9.5348,
		MinTriggerValue: 0.2,
		MaxTriggerValue: 1.3,
	}

	err := wh.Invoke()
	if err != nil {
		t.Errorf("Invoke failed: %s", err.Error())
	}
}
