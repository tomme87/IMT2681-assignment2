package api

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"net/url"
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
)

type Webhook struct {
	ID bson.ObjectId `json:"_id" bson:"_id"`
	WebhookURL string
	BaseCurrency string
	TargetCurrency string
	MinTriggerValue int
	MaxTriggerValue int
}

func (wh *Webhook) Validate() error {
	if wh.ID.Hex() == "" || wh.WebhookURL == "" || wh.BaseCurrency == "" || wh.TargetCurrency == "" ||
		wh.MinTriggerValue < 0 || wh.MaxTriggerValue < 0 {
		return errors.New("missing input")
	}

	if wh.MinTriggerValue > wh.MaxTriggerValue {
		return errors.New("min higher than max")
	}

	u, err := url.Parse(wh.WebhookURL)
	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("only HTTPS or HTTP urls supported")
	}

	return nil
}

func (wh *Webhook) Invoke() error {
	jsonData, err := json.Marshal(wh)
	if err != nil {
		return err
	}
	res, err := http.Post(wh.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("bad response code. got: %d. body: %s", res.StatusCode, body)
	}

	return nil
}