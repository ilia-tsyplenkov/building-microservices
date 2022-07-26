package handlers

import (
	"net/http"

	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
// responses:
//  201: noContentResponse
//  422: errorValidation

// Create handles POST requests to add new product to the data storage
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Debug("Creating product %#v\n", product)
	p.productDB.AddProduct(product)
}
