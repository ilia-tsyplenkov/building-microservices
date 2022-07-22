package handlers

import (
	"net/http"

	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//  201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// UpdateProduct handles PUT requests to update products
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(KeyProduct{}).(data.Product)
	if err := data.UpdateProduct(product); err != nil {
		http.Error(rw, "error updating product: "+err.Error(), http.StatusBadRequest)
		return

	}
}
