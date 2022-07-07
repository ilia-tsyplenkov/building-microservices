package handlers

import (
	"building-microservices/product-api/data"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		matches := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(matches) != 1 {
			msg := fmt.Sprintf("Invalid URL more than one id: %s", r.URL.Path)
			p.l.Println(msg)
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}
		if len(matches[0]) != 2 {
			msg := fmt.Sprintf("Invalid URL, more than one captured id group: %s", r.URL.Path)
			p.l.Println(msg)
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(matches[0][1])
		if err != nil {
			msg := fmt.Sprintf("error converting id to int: %s", matches[0][1])
			p.l.Println(msg)
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}
		p.updateProduct(id, rw, r)
		return
	}
	rw.WriteHeader(http.StatusNotImplemented)

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
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

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	product := data.NewProduct()
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "error marshaling request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	data.AddProduct(product)

}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	product := data.NewProduct()
	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "error marshaling request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	// p.l.Printf("got product to update: %#v\n", product)
	if err := data.UpdateProduct(id, product); err != nil {
		http.Error(rw, "error updating product: "+err.Error(), http.StatusBadRequest)
		return

	}
}
