package main

import (
	"building-microservices/product-api/sdk/client"
	"building-microservices/product-api/sdk/client/products"
	"testing"
)

var cfg = client.DefaultTransportConfig().WithHost("localhost:8080")
var cl = client.NewHTTPClientWithConfig(nil, cfg)

func TestClientListSingleProduct(t *testing.T) {
	params := products.NewListSingleProductParams()
	params.ID = 2
	_, err := cl.Products.ListSingleProduct(params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientListProducts(t *testing.T) {
	params := products.NewListProductsParams()
	_, err := cl.Products.ListProducts(params)
	if err != nil {
		t.Fatal(err)
	}
}
