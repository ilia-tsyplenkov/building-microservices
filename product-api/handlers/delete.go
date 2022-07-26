package handlers

import (
	"fmt"
	"net/http"

	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete the product
// responses:
//  201: noContentResponse
//  404: errorResponse

// Delete handles DELETE requests to delete the product
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := p.getProductID(r)
	p.l.Debug("deleting product", "id", id)
	if err := p.productDB.DeleteProduct(id); err != nil {
		msg := fmt.Sprintf("error deleting product: %s", err.Error())
		p.l.Error("[ERROR] " + msg)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: msg}, rw)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
