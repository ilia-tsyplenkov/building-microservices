package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

const ratesURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewExchangeRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rates: map[string]float64{}}
	err := er.getRates()
	return er, err
}

func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %q", base)
	}
	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %q", dest)
	}
	return dr / br, nil
}
func (e *ExchangeRates) getRates() error {
	resp, err := http.Get(ratesURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected to have 200 status code, but got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	md := Cubes{}
	err = decoder.Decode(&md)
	if err != nil {
		return err
	}

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return fmt.Errorf("error converting rate from str to float: %v", err)
		}
		e.rates[c.Currency] = r
	}
	e.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
