package data

import (
	"fmt"
	"time"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id of the product

	// required: false
	// min: 1
	ID int `json:"id"`
	// the name of the product

	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`
	// the description of the product

	// required: false
	// max length: 10000
	Description string `json:"description"`
	// the price of the product

	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"required,gt=0"`
	// the sku of the product

	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"required,sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

func NewProduct() Product {
	return Product{CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	}
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

type Products []*Product

func GetProducts() Products {
	return productList
}

func GetProduct(id int) (*Product, error) {
	p, _, err := findProduct(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

func getNextId() int {
	last := productList[len(productList)-1]
	return last.ID + 1
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func UpdateProduct(p Product) error {
	oldP, pos, err := findProduct(p.ID)
	if err != nil {
		return err
	}
	p.CreatedOn = oldP.CreatedOn
	productList[pos] = &p
	return nil
}

func DeleteProduct(id int) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	if pos == len(productList)-1 {
		productList = productList[:pos]
	} else {
		productList = append(productList[:pos], productList[pos+1:]...)
	}
	return nil
}
