package api

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

// Fixer holds rates from fixer
type Fixer struct {
	Base string
	Date string
	Rates map[string]float32
}

// NewFixer create the fixer object with data from latest. with optional url.
func NewFixer(url ...string) (*Fixer, error) {
	f := new(Fixer)

	if len(url) == 0 {
		url = append(url, "http://api.fixer.io/latest")
	}

	res, err := http.Get(url[0])
	if err != nil {
		return f, err
	}

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return f, fmt.Errorf("bad response code. got: %d. body: %s", res.StatusCode, body)
	}

	err = json.NewDecoder(res.Body).Decode(f)
	if err != nil {
		return f, err
	}

	return f, nil
}

// GetRate to convert rate from base to target.
func (f *Fixer) GetRate(base string, target string) (float32, error) {
	if _, ok := f.Rates[base]; !ok && f.Base != base {
		return 0, fmt.Errorf("not av valid base: %s", base)
	}

	if _, ok := f.Rates[target]; !ok && f.Base != target{
		return 0, fmt.Errorf("not a valid target: %s", target)
	}

	if base == f.Base {
		return f.Rates[target], nil
	}

	if target == f.Base {
		return 1 / f.Rates[base], nil
	}

	return f.Rates[target] * (1 / f.Rates[base]), nil
}