package data

import (
	"testing"
)

func TestCheckValidation(t *testing.T) {
	p := &Product{Name: "Ilia",
		Price: 1.0,
		SKU:   "aaa-bbb-ccc"}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
