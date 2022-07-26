package data

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	protos "github.com/ilia-tsyplenkov/building-microservices/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id of the product

	// required: false
	// min: 1
	ID int `json:"id"`
	// the name of the product

	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`
	// the description of the product

	// required: false
	// max length: 10000
	Description string `json:"description"`
	// the price of the product

	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`
	// the sku of the product

	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"required,sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

func NewProduct() Product {
	return Product{CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	}
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

type Products []*Product

type ProductsDB struct {
	currency     protos.CurrencyClient
	log          hclog.Logger
	rates        map[string]float64
	m            sync.Mutex
	streamClient protos.Currency_SubscribeRatesClient
}

func NewProductDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{currency: c, log: l, rates: make(map[string]float64)}
	go pb.handleUpdates()
	return pb
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("unable to subscribe for rates", "error", err)
		return
	}
	p.streamClient = sub
	for {
		rr, err := sub.Recv()
		if err != nil {
			p.log.Error("unable to receive rates", "error", err)
			return
		}
		if grpcErr := rr.GetError(); grpcErr != nil {
			p.log.Error("error subscribing for rates", "error", grpcErr)
			continue
		}
		if resp := rr.GetRateResponse(); resp != nil {
			p.log.Info("received updated rate from server", "destination", resp.Destination.String())
			p.m.Lock()
			p.rates[resp.Destination.String()] = resp.Rate
			p.m.Unlock()

		}

	}

}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" || currency == "EUR" {
		return productList, nil
	}
	rate, err := p.getRate("EUR", currency)
	if err != nil {
		p.log.Error("error getting exchange rate", "currency", currency, "error", err)
		return nil, err
	}
	resultProductList := make(Products, 0)
	for _, p := range productList {
		product := *p
		product.Price = product.Price * rate
		resultProductList = append(resultProductList, &product)
	}

	return resultProductList, nil
}

func (p *ProductsDB) GetProduct(id int, currency string) (*Product, error) {
	product, _, err := findProduct(id)
	if err != nil {
		return nil, err
	}
	if currency == "" || currency == "EUR" {
		return product, nil
	}
	rate, err := p.getRate("EUR", currency)
	if err != nil {
		p.log.Error("error getting exchange rate", "currency", currency, "error", err)
		return nil, err
	}
	res := *product
	res.Price = res.Price * rate
	return &res, nil
}

func (p *ProductsDB) AddProduct(product Product) {
	product.ID = getNextId()
	productList = append(productList, &product)
}

func (p *ProductsDB) UpdateProduct(product Product) error {
	oldP, pos, err := findProduct(product.ID)
	if err != nil {
		return err
	}
	product.CreatedOn = oldP.CreatedOn
	productList[pos] = &product
	return nil
}

func (p *ProductsDB) DeleteProduct(id int) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	if pos == len(productList)-1 {
		productList = productList[:pos]
	} else {
		productList = append(productList[:pos], productList[pos+1:]...)
	}
	return nil
}

func (p *ProductsDB) getRate(baseCurrency, dstCurrency string) (float64, error) {
	// if cached, return
	p.m.Lock()
	defer p.m.Unlock()
	if r, ok := p.rates[dstCurrency]; ok {
		return r, nil
	}
	// get exchange rate
	request := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value[baseCurrency]),
		Destination: protos.Currencies(protos.Currencies_value[dstCurrency]),
	}
	// get initial rates
	rr, err := p.currency.GetRate(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*protos.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("Unable to get rate from currency server, destination and base currencies can not be the same, base: %s, dest: %s",
					md.Base.String(),
					md.Destination.String())

			}
			return -1, fmt.Errorf("Unable to get rate from currency server, base: %s, dest: %s",
				md.Base.String(),
				md.Destination.String())
		}
		return -1, err
	}
	p.rates[dstCurrency] = rr.Rate // update cache

	// subscribe for updates
	err = p.streamClient.Send(request)
	if err != nil {
		p.log.Error("unable to subscribe for updates", "destination", dstCurrency, "error", err)
	}
	return rr.Rate, err

}

func getNextId() int {
	last := productList[len(productList)-1]
	return last.ID + 1
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}
