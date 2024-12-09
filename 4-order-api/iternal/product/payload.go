package product

import "github.com/lib/pq"

type ProductCreateRequest struct {
	Name        string         `json:"url" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"  validate:"required"`
}

type ProductUpdateRequest struct {
	Name        string         `json:"url"`
	Description string         `json:"description"`
	Images      pq.StringArray `gorm:"type:text[]" json:"images"`
}
