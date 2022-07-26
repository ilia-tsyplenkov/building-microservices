package handlers

import (
	"context"
	"net/http"

	"github.com/ilia-tsyplenkov/building-microservices/product-api/data"
)

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := data.NewProduct()
		if err := data.FromJSON(&product, r.Body); err != nil {
			msg := "error marshaling request body: " + err.Error()
			p.l.Error("[ERROR]: " + msg)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: msg}, rw)
			return
		}
		if errs := p.v.Validate(&product); errs != nil && len(errs) != 0 {
			msg := "error validate the product: "
			p.l.Error(msg, "errors", errs.Errors())
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})

}
