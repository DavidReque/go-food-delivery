package data

import (
	"context"

	"github.com/DavidReque/go-food-delivery/internal/services/catalogreadservice/internal/products/models"
)

// ProductCacheRepository is the interface for the product cache repository
type ProductCacheRepository interface {
	PutProduct(ctx context.Context, key string, product *models.Product) error
	GetProductById(ctx context.Context, key string) (*models.Product, error)
	DeleteProduct(ctx context.Context, key string) error
	DeleteAllProducts(ctx context.Context) error
}
