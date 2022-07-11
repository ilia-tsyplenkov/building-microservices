package data

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckValidation(t *testing.T) {
	testCases := []struct {
		caseName string
		name     string
		price    float32
		sku      string
		errLen   int
	}{
		{caseName: "MissingNamePriceSKUReturnsErr", errLen: 3},
		{caseName: "WrongPriceReturnErr", name: "test", sku: "a-b-c", price: -0.01, errLen: 1},
		{caseName: "MissingPriceReturnErr", name: "test", sku: "a-b-c", errLen: 1},
		{caseName: "MissingNameReturnErr", price: 1.99, sku: "a-b-c", errLen: 1},
		{caseName: "MissingSKUReturnErr", price: 1.99, name: "test", errLen: 1},
		{caseName: "WrongSKUReturnErr", price: 1.99, name: "test", sku: "wrong-sku", errLen: 1},
		{caseName: "ValidProductDoesNOTReturnErr", price: 1.99, name: "test", sku: "a-b-c", errLen: 0},
	}

	v := NewValidation()
	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			p := Product{
				Name:  tc.name,
				Price: tc.price,
				SKU:   tc.sku,
			}
			err := v.Validate(p)
			assert.Len(t, err, tc.errLen)
		})
	}
}

func TestProductToJSON(t *testing.T) {
	p := Product{
		Name:  "foo",
		Price: 1.99,
		SKU:   "a-b-c",
	}
	buf := bytes.NewBufferString("")
	err := ToJSON(&p, buf)
	assert.NoError(t, err)
}

func TestProductFromJSON(t *testing.T) {
	p := Product{}
	buf := bytes.NewBufferString(`{"name":"test", "price": 1.99, "sku": "a-b-c"}`)
	err := FromJSON(&p, buf)
	assert.NoError(t, err)
	assert.Equal(t, p.Name, "test")
	assert.Equal(t, p.Price, float32(1.99))
	assert.Equal(t, p.SKU, "a-b-c")
}
