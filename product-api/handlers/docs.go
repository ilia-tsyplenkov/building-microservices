// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import "building-microservices/product-api/data"

// A list of products returns in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// Product structure returns in the response
// swagger:response productResponse
type productResponseWrapper struct {
	// Products in the system
	// in: body
	Body data.Product
}

// swagger:parameters deleteProduct listSingleProduct
type productIDParameterWrapper struct {
	// The id for the product to delete from database
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of errors
	// in: body
	Body ValidationError
}

// swagger:response errorResponse
type errorResponseWrapper struct {
	// Error message
	// in: body
	Body GenericError
}

// swagger:parameters updateProduct createProduct
type productParamsWrapper struct {
	// Product data structure to Update or Create.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body data.Product
}
