package data

import (
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewRates(t *testing.T) {
	r, err := NewExchangeRates(hclog.Default())
	if err != nil {
		t.Fatalf("expected nil but got %q", err.Error())
	}
	if r.rates == nil || len(r.rates) == 0 {
		t.Fatalf("expected to have some rates data, but got %v", r.rates)
	}

}
