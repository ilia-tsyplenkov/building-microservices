package handlers

import (
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"

	"github.com/gorilla/mux"
)

type Products struct {
	l         hclog.Logger
	v         *data.Validation
	productDB *data.ProductsDB
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

func NewProducts(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
}
