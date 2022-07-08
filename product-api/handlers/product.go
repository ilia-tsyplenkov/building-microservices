package handlers

import (
	"building-microservices/product-api/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJSON(rw)
	if err != nil {
		var status int
		if err == data.ErrProductNotFound {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		http.Error(rw, err.Error(), status)
		return
	}

}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&product)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("error converting id to int: %s", vars["id"])
		p.l.Println(msg)
		http.Error(rw, msg, http.StatusBadRequest)
		return

	}
	product := r.Context().Value(KeyProduct{}).(data.Product)
	if err := data.UpdateProduct(id, &product); err != nil {
		http.Error(rw, "error updating product: "+err.Error(), http.StatusBadRequest)
		return

	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := data.NewProduct()
		if err := product.FromJSON(r.Body); err != nil {
			msg := "error marshaling request body: " + err.Error()
			p.l.Println("[ERROR]: " + msg)
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}
		if err := product.Validate(); err != nil {
			msg := "error validating product fields: " + err.Error()
			p.l.Println("[ERROR]: " + msg)
			http.Error(rw, msg, http.StatusBadRequest)
			return

		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})

}
