package handlers

import (
	"building-microservices/product-api/data"
	"fmt"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete the product
// responses:
//  201: noContentResponse
//  404: errorResponse

// Delete handles DELETE requests to delete the product
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := p.getProductID(r)
	if err := data.DeleteProduct(id); err != nil {
		msg := fmt.Sprintf("error deleting product: %s", err.Error())
		p.l.Println("[ERROR] " + msg)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: msg}, rw)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
