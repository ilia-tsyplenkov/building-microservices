package handlers

import (
	"building-microservices/product-api/data"
	"fmt"
	"net/http"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//  200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	rw.Header().Add("Content-Type", "application/json")
	err := data.ToJSON(products, rw)
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

// swagger:route GET /products/{id} products listSingleProduct
// Returns the product
// responses:
//  200: productResponse
//  404: errorResponse

// GetProduct returns the product from the data store
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	id := p.getProductID(r)
	rw.Header().Add("Content-Type", "application/json")
	product, err := data.GetProduct(id)
	if err != nil {
		msg := fmt.Sprintf("product not found: product Id - %d", id)
		p.l.Println("[ERROR] " + msg)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: msg}, rw)
		return

	}
	if err := data.ToJSON(product, rw); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

}
