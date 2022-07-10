package handlers

import (
	"building-microservices/product-api/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
	v *data.Validation
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

func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}