package clients

import "context"

type ProductDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Category    int     `json:"category"`
	ImageLink   string  `json:"image_link"`
}

type ProductClient interface {
	GetProduct(ctx context.Context, productID uint) (*ProductDTO, error)
	GetProducts(ctx context.Context, productIDs []uint) ([]*ProductDTO, error)
	ValidateProduct(ctx context.Context, productID uint) (bool, error)
}
