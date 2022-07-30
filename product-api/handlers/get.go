package handlers

import (
	"net/http"

	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//  200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get all records")
	rw.Header().Add("Content-Type", "application/json")
	currency := r.URL.Query().Get("currency")
	if currency != "" {
		p.l.Debug("non default currency is provided. all prices would be shown in new one.", "currency", currency)
	} else {
		p.l.Debug("no currency provided. all prices would be shown in EUR.")
	}
	products, err := p.productDB.GetProducts(currency)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return

	}
	err = data.ToJSON(products, rw)
	if err != nil {
		p.l.Error("Unable to serialize product", "error", err)
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
	p.l.Debug("Get record", "id", id)
	rw.Header().Add("Content-Type", "application/json")
	currency := r.URL.Query().Get("currency")
	if currency != "" {
		p.l.Debug("non default currency is provided. all prices would be shown in new one.", "currency", currency)
	} else {
		p.l.Debug("no currency provided. all prices would be shown in EUR.")
	}
	product, err := p.productDB.GetProduct(id, currency)
	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Error("unable to find product", "id", id, "error", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("unable to fetch product", "id", id, "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
	if err := data.ToJSON(product, rw); err != nil {
		p.l.Error("Unable to serialize product", "id", id, "error", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

}
