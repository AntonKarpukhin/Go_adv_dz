package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
}

func NewProduct(product *ProductCreateRequest) *Product {
	return &Product{
		Name:        product.Name,
		Description: product.Description,
		Images:      product.Images,
	}
}
