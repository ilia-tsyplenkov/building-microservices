package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"

	"github.com/gorilla/mux"
	protos "github.com/ilia-tsyplenkov/building-microservices/currency/protos/currency"
)

type Products struct {
	l  *log.Logger
	v  *data.Validation
	cc protos.CurrencyClient
}

func (p *Products) getProductID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)

	}
	return id
}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

func NewProducts(l *log.Logger, v *data.Validation, cc protos.CurrencyClient) *Products {
	return &Products{l, v, cc}
}
