package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func NewProduct() *Product {
	p := &Product{CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	}
	return p
}
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

var ErrProductNotFound = fmt.Errorf("Product not found")

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

func (p *Products) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

func GetProducts() Products {
	return productList
}

func UpdateProducts(update Products) {
	for _, u := range update {
		updated := false
		for i, p := range productList {
			if u.ID == p.ID {
				tmp := p
				p = u
				p.CreatedOn = tmp.CreatedOn
				productList[i] = p
				updated = true
			}
		}
		if !updated {
			productList = append(productList, u)
		}
	}
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

func UpdateProduct(id int, p *Product) error {
	oldP, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.CreatedOn = oldP.CreatedOn
	p.ID = id
	productList[pos] = p
	return nil
}
